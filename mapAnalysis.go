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

const maxLenght int16 = 200
const maxWidth int16 = 200

type pdfMap struct {
	points    []readcoordinates.Point
	propotion float64
	width     int16
	lenght    int16
}

func newPoint(points []readcoordinates.Point, propotion float64) *pdfMap {
	m := pdfMap{propotion: propotion, points: points, lenght: maxLenght}
	return &m
}
func (m *pdfMap) defineWidth() {
	m.width = int16(float64(m.lenght)/m.propotion) + 2
	if m.width > maxWidth {
		m.width = int16(float64(m.width)*m.propotion) + 1

	} else {
		return
	}

}
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
