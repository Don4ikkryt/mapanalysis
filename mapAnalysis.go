package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/Don4ikkryt/readcoordinates"
	"github.com/jung-kurt/gofpdf"
)

var (
	sourceFolder             string
	mapFolder                string
	filteredFolder           string
	maxDistanceBetweenPoints int64
)

const maxLenght int16 = 150
const maxWidth int16 = 150

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
	i := 0
	for _, value := range m.centreInPoint.Latitude {
		value /= 2
		value += readcoordinates.Southest.Latitude[i]
		i++
	}

	m.centreInPoint.Longtitude = readcoordinates.CoordinateDiffernce(readcoordinates.Westest.Longtitude, readcoordinates.Eastest.Longtitude)
	j := 0
	for _, value := range m.centreInPoint.Longtitude {
		value /= 2
		value += readcoordinates.Westest.Longtitude[i]
		j++
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

	pdfFile := openPDFFile()
	createRectangle(pdfFile, map1.lenght, map1.width)
	closePDFFile(pdfFile)

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

func closePDFFile(f *gofpdf.Fpdf) {
	var i int = 1
	var filename string
	for {
		filename = mapFolder + "\\map" + strconv.Itoa(i)
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			break
		}
		i++
	}

	err := f.OutputFileAndClose(filename)
	if err != nil {
		fmt.Println(1)
		fmt.Println("err")
	}
}
func openPDFFile() (f *gofpdf.Fpdf) {
	f = gofpdf.New("P", "mm", "A4", "")
	f.AddPage()
	return
}
func createRectangle(f *gofpdf.Fpdf, lenght int16, width int16) {
	f.SetFillColor(0, 153, 255)
	f.Rect(25, 25, float64(lenght), float64(width), "F")
}
func (m *pdfMap) ConvertFromPointsToMMPoint(point readcoordinates.Point) (newPoint mmPoint) {
	DistanceInMetersLatitude := readcoordinates.ConvertFromCoordinatesToMeterLatitude(readcoordinates.CoordinateDiffernce(point.Latitude, m.centreInPoint.Latitude))
	DistanceInMetersLongtitude := readcoordinates.ConvertFromCoordinatesToMeterLongtitude(readcoordinates.CoordinateDiffernce(point.Longtitude, m.centreInPoint.Longtitude), &point, &m.centreInPoint)

}
