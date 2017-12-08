package main

import (
	"flag"
	"fmt"
	"github.com/fatih/color"
	"github.com/qrtp/gpxgo/gpx"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	AggregateHistogram = "AGGREGATE_DATA"
	GroupByMonth       = "month"
	GroupByWeek        = "week"
)

func main() {

	// Parse flags
	fileList := flag.String("files", "default.gpx", "comma separated list of paths to GPX files")
	dirPath := flag.String("directory", "", "path to directory to search for GPX files")
	zoneList := flag.String("zones", "150", "comma separated list of heart rate zone thresholds")
	groupBy := flag.String("groupBy", "month", "can specify either month or week, default=month")
	flag.Parse()

	// Validate group flag
	if *groupBy != GroupByWeek && *groupBy != GroupByMonth {
		fmt.Printf("invalid grouping value: %s\n", *groupBy)
		os.Exit(1)
	}

	// Create aggregate zone buckets
	zoneThresholds := strings.Split(*zoneList, ",")

	// Heart rate zones by group, plus a special aggregate histogram
	// to summarize the whole data set.
	histograms := make(map[string]*Histogram)
	histograms[AggregateHistogram] = newHistogram("Heart Rate Zone Summary", zoneThresholds)

	// Search files if dirpath provided
	if dirPath != nil && *dirPath != "" {
		files, err := filepath.Glob(fmt.Sprintf("%s", *dirPath))
		if err != nil {
			fmt.Printf("unable to list directory %s: %s\n", *dirPath, err.Error())
			os.Exit(2)
		}
		if len(files) == 0 {
			fmt.Printf("unable to locate any GPX files in directory: %s\n", *dirPath)
			os.Exit(2)
		}
		joinedFiles := strings.Join(files, ",")
		fileList = &joinedFiles
	}

	// Iterate files
	for _, pathToFile := range strings.Split(*fileList, ",") {

		// Is the file a GPX file?
		if filepath.Ext(strings.ToLower(pathToFile)) != ".gpx" {
			continue
		}

		// Read provided file
		gpxBytes, err := ioutil.ReadFile(pathToFile)
		if err != nil {
			fmt.Printf("unable to read GPX file %s: %s\n", pathToFile, err.Error())
			os.Exit(3)
		}
		gpxFile, err := gpx.ParseBytes(gpxBytes)
		if err != nil {
			fmt.Printf("unable to parse GPX file %s: %s\n", pathToFile, err.Error())
			os.Exit(4)
		}

		// Determine group key
		groupKey := fmt.Sprintf("%d-%s", gpxFile.Time.Year(), fmt.Sprintf("%02d", gpxFile.Time.Month()))
		if *groupBy == GroupByWeek {
			year, weekOfYear := gpxFile.Time.ISOWeek()
			groupKey = fmt.Sprintf("%d week %s", year, fmt.Sprintf("%02d", weekOfYear))
		}

		// Initialize group histogram if this is first data point
		if _, ok := histograms[groupKey]; !ok {
			histograms[groupKey] = newHistogram(fmt.Sprintf("%s Heart rate zones", groupKey), zoneThresholds)
		}

		// Dump GPX heart rates into zone histograms
		for _, track := range gpxFile.Tracks {
			for _, segment := range track.Segments {
				for i := 0; i < len(segment.Points)-1; i++ {
					currPoint := segment.Points[i]
					nextPoint := segment.Points[i+1]
					secondsElapsed := currPoint.TimeDiff(&nextPoint)
					histograms[AggregateHistogram].AddHeartRate(currPoint.HeartRate, secondsElapsed)
					histograms[groupKey].AddHeartRate(currPoint.HeartRate, secondsElapsed)
				}
			}
		}
	}

	// Render the results
	numHistograms := len(histograms)
	if numHistograms > 2 {
		groupKeys := make([]string, 0, len(histograms))
		for key := range histograms {
			groupKeys = append(groupKeys, key)
		}
		sort.Strings(groupKeys)
		for _, key := range groupKeys {
			histograms[key].Print()
		}
	} else {
		histograms[AggregateHistogram].Print()
	}
}

type Bucket struct {
	Count        int
	ThresholdHR  int
	TotalSeconds float64
}

type Histogram struct {
	Buckets          []Bucket
	Name             string
	MaxEncounteredHR int
	TotalSeconds     float64
	TotalHR          float64
}

func newHistogram(name string, zoneThresholds []string) *Histogram {
	histogram := Histogram{
		Name:    name,
		Buckets: make([]Bucket, len(zoneThresholds)+1),
	}
	for i, threshold := range zoneThresholds {
		hr, err := strconv.Atoi(threshold)
		if err != nil {
			fmt.Printf("unable to create zone %s\n", threshold)
			os.Exit(1)
		}
		histogram.Buckets[i+1].ThresholdHR = hr
	}
	return &histogram
}

func (h *Histogram) Print() {
	total := 0.0
	for _, bucket := range h.Buckets {
		total += bucket.TotalSeconds
	}

	// Render the histogram
	fmtGreen := color.New(color.FgGreen).Add(color.Bold)
	fmtGreen.Println(h.Name)
	for _, bucket := range h.Buckets {
		pct := 100 * bucket.TotalSeconds / total
		duration := time.Duration(bucket.TotalSeconds) * time.Second
		fmt.Printf("> %d \t %.1f%% \t [%s]\n", bucket.ThresholdHR, pct, duration.String())
	}
	fmt.Println("")
	fmt.Printf("Max: %d\n", h.MaxEncounteredHR)
	fmt.Printf("Avg: %.0f\n", h.TotalHR/h.TotalSeconds)
	fmt.Println("")
}

func (h *Histogram) AddHeartRate(hr int, elapsed float64) {

	// Record cumulative statistics
	if hr > h.MaxEncounteredHR {
		h.MaxEncounteredHR = hr
	}
	h.TotalSeconds += elapsed
	h.TotalHR += float64(hr) * elapsed

	// Choose and add to bucket
	numBuckets := len(h.Buckets)
	for i := 0; i < numBuckets-1; i++ {
		if hr < h.Buckets[i+1].ThresholdHR {
			bucket := &h.Buckets[i]
			bucket.AddToBucket(elapsed)
			return
		}
	}
	bucket := &h.Buckets[numBuckets-1]
	bucket.AddToBucket(elapsed)
}

func (b *Bucket) AddToBucket(elapsed float64) {
	b.Count++
	b.TotalSeconds += elapsed
}
