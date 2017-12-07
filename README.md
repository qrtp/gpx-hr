# GPX Heart Rate Zone Aggregator

This tool is a **very** rough swag at analyzing heart rate zones across workouts. Most workout tracking platforms will
provide per-workout zone analysis but an aggregate view is not available at this time. Please provide feedback if you
have any. Enjoy.

## Download Strava activities
Instructions to download all activities from an account can be found [here](https://support.strava.com/hc/en-us/articles/216918437-Exporting-your-Data-and-Bulk-Export). Follow
these instructions and extract the downloaded zip file onto your computer.

## Clone this project

From the command line:
```
git clone git@github.com:qrtp/gpx-hr.git
cd gpx-hr
```

## Analyze heart rate zones across multiple workouts

The `directory` and `zones` flags both need to be specified. The `directory` is the path where you extracted all the GPX
files downloaded from Strava. The `zones` flag is a comma separated list of heart rate zones to analyze. This of course
is different for every athlete. I wrote this tool to determine how close I am to 80/20 running, and as you can tell below
my pacing needs a **LOT** of work based on a 150 bpm threshold.

```
go run gpx-hr.go --directory "/path/to/strava/activities/*.gpx" --zones 150
```

Results in the following output:

```
2017-08 Heart rate zones
> 0 	 11.8% 	 [38m31s]
> 150 	 88.2% 	 [4h49m13s]
2017-09 Heart rate zones
> 0 	 12.3% 	 [1h4m45s]
> 150 	 87.7% 	 [7h43m12s]
2017-10 Heart rate zones
> 0 	 19.3% 	 [2h43m9s]
> 150 	 80.7% 	 [11h20m46s]
2017-11 Heart rate zones
> 0 	 14.4% 	 [1h56m11s]
> 150 	 85.6% 	 [11h28m53s]
2017-12 Heart rate zones
> 0 	 31.1% 	 [1h5m26s]
> 150 	 68.9% 	 [2h25m3s]

Aggregate Heart Rate Zones
> 0 	 16.5% 	 [7h28m2s]
> 150 	 83.5% 	 [37h47m7s]
```
