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
go run gpx-hr.go --directory "/path/to/strava/activities/*.gpx" --zones 157,188
```

Results in the following output:

```
2017-08 Heart rate zones
> 0 	 17.3% 	 [56m45s]
> 157 	 81.0% 	 [4h25m28s]
> 188 	 1.7% 	 [5m31s]

Max: 198
Avg: 168

2017-09 Heart rate zones
> 0 	 24.0% 	 [2h6m46s]
> 157 	 75.4% 	 [6h38m19s]
> 188 	 0.5% 	 [2m52s]

Max: 202
Avg: 162

2017-10 Heart rate zones
> 0 	 37.4% 	 [5h15m28s]
> 157 	 61.3% 	 [8h37m3s]
> 188 	 1.4% 	 [11m24s]

Max: 208
Avg: 158

2017-11 Heart rate zones
> 0 	 33.2% 	 [4h26m53s]
> 157 	 66.2% 	 [8h53m15s]
> 188 	 0.6% 	 [4m56s]

Max: 200
Avg: 160

2017-12 Heart rate zones
> 0 	 60.0% 	 [2h43m1s]
> 157 	 37.9% 	 [1h43m1s]
> 188 	 2.1% 	 [5m46s]

Max: 206
Avg: 151

Heart Rate Zone Summary
> 0 	 33.5% 	 [15h28m53s]
> 157 	 65.4% 	 [30h17m6s]
> 188 	 1.1% 	 [30m29s]

Max: 208
Avg: 160
```

If desired, there is also an option to group the data by week instead of month. Use the `--groupBy=week` flag 
to enable this output format.