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

func main() {

	// Parse flags
	fileList := flag.String("files", "default.gpx", "comma separated list of paths to GPX files")
	dirPath := flag.String("directory", "", "path to directory to search for GPX files")
	zoneList := flag.String("zones", "150", "comma separated list of heart rate zone thresholds")
	flag.Parse()

	// Create aggregate zone buckets
	zoneThresholds := strings.Split(*zoneList, ",")
	aggregateData := &Histogram{
		Name:    "Aggregate Heart Rate Zones",
		Buckets: initBuckets(zoneThresholds),
	}

	// Month zone buckets
	monthBuckets := make(map[string]Histogram)

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

		// Month data
		monthKey := fmt.Sprintf("%d-%s", gpxFile.Time.Year(), fmt.Sprintf("%02d", gpxFile.Time.Month()))
		if _, ok := monthBuckets[monthKey]; !ok {
			monthBuckets[monthKey] = Histogram{
				Name:    fmt.Sprintf("%s Heart rate zones", monthKey),
				Buckets: initBuckets(zoneThresholds),
			}
		}
		monthHs := monthBuckets[monthKey]
		monthHsPtr := &monthHs

		// Dump GPX heart rates into buckets
		for _, track := range gpxFile.Tracks {
			for _, segment := range track.Segments {
				for i := 0; i < len(segment.Points)-1; i++ {
					currPoint := segment.Points[i]
					nextPoint := segment.Points[i+1]
					secondsElapsed := currPoint.TimeDiff(&nextPoint)
					aggregateData.AddHeartRate(currPoint.HeartRate, secondsElapsed)
					monthHsPtr.AddHeartRate(currPoint.HeartRate, secondsElapsed)
				}
			}
		}
	}

	// Render the results
	monthBucketKeys := make([]string, 0, len(monthBuckets))
	for key := range monthBuckets {
		monthBucketKeys = append(monthBucketKeys, key)
	}
	sort.Strings(monthBucketKeys)
	for _, key := range monthBucketKeys {
		monthBuckets[key].Print()
	}
	fmt.Println("")
	aggregateData.Print()
}

type Bucket struct {
	HeartRate int
	Duration  float64
	Count     int
}

type Histogram struct {
	Name    string
	Buckets []Bucket
}

func initBuckets(zoneThresholds []string) []Bucket {
	buckets := make([]Bucket, len(zoneThresholds)+1)
	for i, threshold := range zoneThresholds {
		hr, err := strconv.Atoi(threshold)
		if err != nil {
			fmt.Printf("unable to create zone %s\n", threshold)
			os.Exit(1)
		}
		buckets[i+1].HeartRate = hr
	}
	return buckets
}

func (h Histogram) Print() {
	total := 0.0
	for _, bucket := range h.Buckets {
		total += bucket.Duration
	}

	// Render the histogram
	fmtGreen := color.New(color.FgGreen).Add(color.Bold)
	fmtGreen.Println(h.Name)
	for _, bucket := range h.Buckets {
		pct := 100 * bucket.Duration / total
		duration := time.Duration(bucket.Duration) * time.Second
		fmt.Printf("> %d \t %.1f%% \t [%s]\n", bucket.HeartRate, pct, duration.String())
	}
}

func (h *Histogram) AddHeartRate(hr int, elapsed float64) {
	numBuckets := len(h.Buckets)
	for i := 0; i < numBuckets-1; i++ {
		if hr < h.Buckets[i+1].HeartRate {
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
	b.Duration += elapsed
}
