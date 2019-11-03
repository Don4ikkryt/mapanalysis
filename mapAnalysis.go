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

const maxLenght int16 = 205
const maxWidth int16 = 205

type pdfMap struct {
	points    []readcoordinates.Point
	propotion float64
	width     int16
	lenght    int16
}

func newMap(points []readcoordinates.Point, propotion float64) *pdfMap {
	m := pdfMap{propotion: propotion, points: points, lenght: maxLenght}
	return &m
}
func (m *pdfMap) defineWidthAndLenght() {
	m.width = int16(float64(m.lenght-5)/m.propotion) + 5
	if m.width > maxWidth {
		m.width = int16(float64(m.width)*m.propotion) + 5
		m.lenght = int16(float64(m.width-5)*m.propotion) + 5
		m.lenght = toOdd(m.lenght)
		m.width = toOdd(m.width)
		return
	} else {
		return
	}

}
func main() {
	parseFlags()
	points, propotion := readcoordinates.GetCoordinatesAndPropotion(sourceFolder, filteredFolder)

	map1 := newMap(points, propotion)
	map1.defineWidthAndLenght()
	fmt.Println(map1.propotion)
	fmt.Println(map1.lenght)
	fmt.Println(map1.width)
	readcoordinates.
}
func parseFlags() {
	flag.StringVar(&sourceFolder, "source_folder", "", "Path to the folder with photos")
	flag.StringVar(&filteredFolder, "filtered_folder", "", "Path to the folder with filtered (unsupported format/no exif data) files")
	flag.Int64Var(&maxDistanceBetweenPoints, "max_distance_between_points", 0, "Maximum allowed distance between two points, where two photoes were taken")

	flag.Parse()
}
func toOdd(number int16) int16 {
	if number%2 == 1 {
		return number
	} else {
		return number + 1
	}
}
