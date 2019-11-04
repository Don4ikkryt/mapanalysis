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

type mmPoint struct {
	x int16
	y int16
}
type pdfMap struct {
	points        []readcoordinates.Point
	propotion     float64
	width         int16
	lenght        int16
	scale         float64
	centreInMM    mmPoint
	centreInPoint readcoordinates.Point
}

func newMap(points []readcoordinates.Point, propotion float64) *pdfMap {
	m := pdfMap{propotion: propotion, points: points, lenght: maxLenght}
	return &m
}
func (m *pdfMap) defineCentreInMM() {
	m.centreInMM.x = int16(m.lenght + 1/2)
	m.centreInMM.y = int16(m.width + 1/2)
}
func (m *pdfMap) defineCentreInPoint() {
	m.centreInPoint.Latitude = readcoordinates.CoordinateDiffernce(readcoordinates.Northest.Latitude, readcoordinates.Southest.Latitude)
	for _, value := range m.centreInPoint.Latitude {
		value /= 2
	}

	m.centreInPoint.Longtitude = readcoordinates.CoordinateDiffernce(readcoordinates.Westest.Longtitude, readcoordinates.Eastest.Longtitude)
	for _, value := range m.centreInPoint.Longtitude {
		value /= 2
	}
}
func (m *pdfMap) defineScale() {
	m.scale = readcoordinates.ConvertFromCoordinatesToMeterLatitude(readcoordinates.CoordinateDiffernce(readcoordinates.Northest.Latitude, readcoordinates.Southest.Latitude)) / float64(m.width)
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
	map1.defineCentreInMM()
	map1.defineCentreInPoint()
	map1.defineScale()
	

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
