package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Don4ikkryt/readcoordinates"
	"github.com/jung-kurt/gofpdf"
)

var (
	sourceFolder             string
	mapFolder                string
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

func newMap(points []readcoordinates.Point, propotion float64) *pdfMap {
	m := pdfMap{propotion: propotion, points: points, lenght: maxLenght}
	return &m
}
func (m *pdfMap) defineWidthAndLenght() {
	m.width = int16(float64(m.lenght) / m.propotion)
	if m.width > maxWidth {
		m.width = int16(float64(m.width) * m.propotion)
		m.lenght = int16(float64(m.width) * m.propotion)

		return
	}
	m.lenght = toOdd(m.lenght)
	m.width = toOdd(m.width)
	return

}
func main() {
	parseFlags()
	points, propotion := readcoordinates.GetCoordinatesAndPropotion(sourceFolder, filteredFolder)

	map1 := newMap(points, propotion)
	map1.defineWidthAndLenght()

}
func parseFlags() {
	flag.StringVar(&sourceFolder, "source_folder", "", "Path to the folder with photos")
	flag.StringVar(&filteredFolder, "filtered_folder", "", "Path to the folder with filtered (unsupported format/no exif data) files")
	flag.Int64Var(&maxDistanceBetweenPoints, "max_distance_between_points", 0, "Maximum allowed distance between two points, where two photoes were taken")
	flag.StringVar(&mapFolder, "map_folder", "", "Folder where PDF file will be created")

	flag.Parse()
}
func toOdd(number int16) int16 {
	if number%2 == 1 {
		return number
	} else {
		return number + 1
	}
}

func createPDFFile(f *gofpdf.Fpdf) {
	var i int = 1
	var filename string = "map"
	for {
		if _, err := os.Stat(mapFolder + "\\" + filename + string(i)); os.IsNotExist(err) {
			break
		}
		i++
	}
	err := f.OutputFileAndClose(mapFolder + "\\" + filename + string(i))
	if err != nil {
		fmt.Println("err")
	}
}
