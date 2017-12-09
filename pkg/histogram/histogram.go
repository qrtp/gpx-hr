package histogram

import (
	"fmt"
	"github.com/fatih/color"
	"os"
	"strconv"
	"time"
)

type Histogram struct {
	Buckets          []Bucket
	Name             string
	MaxEncounteredHR int
	TotalSeconds     float64
	TotalHR          float64
}

type Bucket struct {
	Count        int
	ThresholdHR  int
	TotalSeconds float64
}

func NewHistogram(name string, zoneThresholds []string) *Histogram {
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
