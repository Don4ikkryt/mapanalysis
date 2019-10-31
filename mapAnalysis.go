package main

import (
	"flag"
)

var (
	sourceFolder             string
	filteredFolder           string
	maxDistanceBetweenPoints int64
)

func main() {
	parseFlags()

}
func parseFlags() {
	flag.StringVar(&sourceFolder, "source_folder", "", "Path to the folder with photos")
	flag.StringVar(&filteredFolder, "filtered_folder", "", "Path to the folder with filtered (unsupported format/no exif data) files")
	flag.Int64Var(&maxDistanceBetweenPoints, "max_distance_between_points", 0, "Maximum allowed distance between two points, where two photoes were taken")

	flag.Parse()
}
