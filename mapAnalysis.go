package main

import (
	"flag"
	"fmt"

	"github.com/Don4ikkryt/readcoordinates"
)

var (
	sourceFolder             string
	filteredFolder           string
	maxDistanceBetweenPoints int64
)

func main() {
	parseFlags()
	points, propotion := readcoordinates.GetCoordinatesAndPropotion(sourceFolder, filteredFolder)

	fmt.Println(propotion)
	for _, value := range points {
		fmt.Println(value)
	}
}
func parseFlags() {
	flag.StringVar(&sourceFolder, "source_folder", "", "Path to the folder with photos")
	flag.StringVar(&filteredFolder, "filtered_folder", "", "Path to the folder with filtered (unsupported format/no exif data) files")
	flag.Int64Var(&maxDistanceBetweenPoints, "max_distance_between_points", 0, "Maximum allowed distance between two points, where two photoes were taken")

	flag.Parse()
}
