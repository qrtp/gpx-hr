package main

import (
	"flag"
	"fmt"
	"github.com/qrtp/gpx-hr/pkg/histogram"
	"github.com/qrtp/gpxgo/gpx"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const (
	AggregateHistogram = "AGGREGATE_DATA"
	GroupByMonth       = "month"
	GroupByWeek        = "week"
)

func main() {

	// Parse flags
	fileList := flag.String("files", "", "comma separated list of paths to GPX files")
	dirPath := flag.String("directory", "", "path to directory to search for GPX files")
	zoneList := flag.String("zones", "157", "comma separated list of heart rate zone thresholds")
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
	histograms := make(map[string]*histogram.Histogram)
	histograms[AggregateHistogram] = histogram.NewHistogram("Heart Rate Zone Summary", zoneThresholds)

	// Default to CWD if no files specified
	if *dirPath == "" && *fileList == "" {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Printf("unable to determine current working directory: %s\n", err.Error())
			os.Exit(1)
		}
		cwd = fmt.Sprintf("%s/*.gpx", cwd)
		dirPath = &cwd
	}

	// Search files if dirpath provided
	if *dirPath != "" {
		files, err := filepath.Glob(fmt.Sprintf("%s", *dirPath))
		if err != nil {
			fmt.Printf("unable to list directory %s: %s\n", *dirPath, err.Error())
			os.Exit(2)
		}
		if len(files) == 0 {
			fmt.Printf("unable to locate any GPX files in directory: %s\n", *dirPath)
			os.Exit(3)
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
			os.Exit(4)
		}
		gpxFile, err := gpx.ParseBytes(gpxBytes)
		if err != nil {
			fmt.Printf("unable to parse GPX file %s: %s\n", pathToFile, err.Error())
			os.Exit(5)
		}

		// Determine group key
		groupKey := fmt.Sprintf("%d-%s", gpxFile.Time.Year(), fmt.Sprintf("%02d", gpxFile.Time.Month()))
		if *groupBy == GroupByWeek {
			year, weekOfYear := gpxFile.Time.ISOWeek()
			groupKey = fmt.Sprintf("%d week %s", year, fmt.Sprintf("%02d", weekOfYear))
		}

		// Initialize group histogram if this is first data point
		if _, ok := histograms[groupKey]; !ok {
			histograms[groupKey] = histogram.NewHistogram(fmt.Sprintf("%s Heart rate zones", groupKey), zoneThresholds)
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
