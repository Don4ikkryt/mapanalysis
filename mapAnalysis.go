package main

import(
	"fmt"
	"flag"

	"github.com/jung-kurt/gofpdf"
)

var (
	sourceFolder   string
	filteredFolder string
	maxDistanceBetweenPoints int16
)

func main(){
	parseFlags()


}
func parseFlags() {
	flag.StringVar(&sourceFolder, "source_folder", "", "Path to the folder with photos")
	flag.StringVar(&filteredFolder, "filtered_folder", "", "Path to the folder with filtered (unsupported format/no exif data) files")
	flag.Int(&maxDistanceBetweenPoints, "max_distance_between_points", 0, "Maximum allowed distance between two points, where two photoes were taken")
	
	flag.Parse()
}