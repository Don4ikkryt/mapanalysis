package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"strconv"

	"github.com/Don4ikkryt/readcoordinates"
	"github.com/jung-kurt/gofpdf"
)

var (
	sourceFolder             string
	mapFolder                string
	filteredFolder           string
	maxDistanceBetweenPoints float64
	amountOfPhotoesInRadius  int64
)

const (
	maxLenght               int16   = 150
	maxWidth                int16   = 150
	lenghtOfEquatorInMeters float64 = 40000000
	degreesInCircle         float64 = 360
	rectIndentionX          float64 = 25
	rectIndentionY          float64 = 25
	basicPointR             int     = 0
	basicPointG             int     = 0
	basicPointB             int     = 0
	NotValidPointR          int     = 255
	NotValidPointG          int     = 0
	NotValidPointB          int     = 0
	colorOfRectangleR       int     = 0
	colorOfRectangleG       int     = 153
	colorOfRectangleB       int     = 255
	radiusOfPoint           float64 = 0.5
)

type mmPoint struct {
	x float64
	y float64
	r int
	g int
	b int
}

func (p *mmPoint) setColor(r, g, b int) {
	p.r = r
	p.g = g
	p.b = b
}

type pdfMap struct {
	points        []readcoordinates.Point
	propotion     float64
	width         int16
	lenght        int16
	scale         float64
	centreInMM    mmPoint
	centreInPoint readcoordinates.Point
	pdfFile       *gofpdf.Fpdf
}

func (m *pdfMap) setWidthAndLenght() {
	m.lenght = maxLenght
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
func (m *pdfMap) setScale() {

	m.scale = float64(m.width) / convertLatitudeToMeters(float64(readcoordinates.PointWithBiggestLatitude.Latitude-readcoordinates.PointWithLeastLatitude.Latitude))
}
func toOdd(number int16) int16 {
	if number%2 == 1 {
		return number
	} else {
		return number + 1
	}
}
func (m *pdfMap) setCentreInPoint() {
	m.centreInPoint.Latitude = (readcoordinates.PointWithBiggestLatitude.Latitude-readcoordinates.PointWithLeastLatitude.Latitude)/2 + readcoordinates.PointWithLeastLatitude.Latitude
	m.centreInPoint.Longtitude = (readcoordinates.PointWithBiggestLongtitude.Longtitude-readcoordinates.PointWithLeastLongtitude.Longtitude)/2 + readcoordinates.PointWithLeastLongtitude.Longtitude
}
func (m *pdfMap) setCentreInMM() {
	m.centreInMM.x = float64(int16(float64(m.lenght)/2 + 1 + rectIndentionX))
	m.centreInMM.y = float64(int16(float64(m.width)/2 + 1 + rectIndentionY))
}
func (m *pdfMap) drawPoint(point *mmPoint) {
	m.pdfFile.SetFillColor(point.r, point.g, point.b)
	m.pdfFile.SetLineWidth(0.02)
	m.pdfFile.Circle(point.x, point.y, radiusOfPoint, "F")
}
func (m *pdfMap) setMap() {

	m.setWidthAndLenght()
	m.setCentreInPoint()
	m.setCentreInMM()
	m.setScale()
}
func (m *pdfMap) DiffBetweenCentreAndPoint(point *readcoordinates.Point) *mmPoint {
	var newMMPoint mmPoint
	newMMPoint.setColor(basicPointR, basicPointG, basicPointB)
	newMMPoint.x = convertLongtitudeToMeters(float64(point.Longtitude-m.centreInPoint.Longtitude), point, &m.centreInPoint)*m.scale + m.centreInMM.x
	newMMPoint.y = convertLatitudeToMeters(float64(point.Latitude-m.centreInPoint.Latitude))*m.scale + m.centreInMM.y

	return &newMMPoint
}

func (m *pdfMap) drawAllPoints(maxDistance float64) {

	for _, value := range m.points {
		fmt.Println(value.Filename)
		point := m.DiffBetweenCentreAndPoint(&value)
		distance := m.distanceBetweenPoints(&value)
		ifNeighbours(point, distance, maxDistance, amountOfPhotoesInRadius)
		m.drawPoint(point)

	}
}
func main() {
	parseFlags()
	points := readcoordinates.GetPoints(sourceFolder, filteredFolder)
	map1 := newMap(points, getPropotion())
	map1.setMap()

	map1.pdfFile = openPDFFile()

	createRectangle(map1.pdfFile, map1.lenght, map1.width)
	map1.drawAllPoints(maxDistanceBetweenPoints)
	closePDFFile(map1.pdfFile)

}
func parseFlags() {
	flag.StringVar(&sourceFolder, "source_folder", "", "Path to the folder with photos")
	flag.StringVar(&filteredFolder, "filtered_folder", "", "Path to the folder with filtered (unsupported format/no exif data) files")
	flag.Float64Var(&maxDistanceBetweenPoints, "max_distance_between_points", 0, "Maximum allowed distance between two points, where two photoes were taken")
	flag.StringVar(&mapFolder, "map_folder", "", "Folder where PDF file will be created")
	flag.Int64Var(&amountOfPhotoesInRadius, "amount_of_photoes_in_radius", 0, "Amount of photoes in the radius to validate it")

	flag.Parse()
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
	f.SetFillColor(colorOfRectangleR, colorOfRectangleG, colorOfRectangleB)
	f.Rect(rectIndentionX, rectIndentionY, float64(lenght)+5, float64(width)+5, "F")
}
func convertLatitudeToMeters(coordinates float64) (meters float64) {
	meters = coordinates * lenghtOfEquatorInMeters / degreesInCircle
	return

}
func convertFromDegreeToRadian(coordinates float64) (radian float64) {
	radian = coordinates * math.Pi / 180
	return
}
func convertLongtitudeToMeters(coordinates float64, point1 *readcoordinates.Point, point2 *readcoordinates.Point) (meters float64) {
	var lesserLatitude float64
	if point1.Latitude < point2.Latitude {
		lesserLatitude = float64(point1.Latitude)
	} else {
		lesserLatitude = float64(point2.Latitude)
	}

	meters = float64(coordinates) * math.Cos(convertFromDegreeToRadian(lesserLatitude)) * lenghtOfEquatorInMeters / degreesInCircle
	return
}

func newMap(points []readcoordinates.Point, propotion float64) *pdfMap {
	m := pdfMap{points: points, propotion: propotion}
	return &m
}
func getPropotion() (propotion float64) {
	longtitudeDiff := readcoordinates.PointWithBiggestLongtitude.Longtitude - readcoordinates.PointWithLeastLongtitude.Longtitude
	lenght := convertLongtitudeToMeters(float64(longtitudeDiff), &readcoordinates.PointWithBiggestLongtitude, &readcoordinates.PointWithLeastLongtitude)

	latitudeDiff := readcoordinates.PointWithBiggestLatitude.Latitude - readcoordinates.PointWithLeastLatitude.Latitude
	width := convertLatitudeToMeters(float64(latitudeDiff))

	propotion = lenght / width
	return
}
func (m *pdfMap) distanceBetweenPoints(point1 *readcoordinates.Point) []float64 {
	pointDistance := make([]float64, 0)

	for _, point2 := range m.points {

		if point1.Filename != point2.Filename {
			longtitudeDiff := convertLongtitudeToMeters(float64(point1.Longtitude-point2.Longtitude), point1, &point2)
			latitudeDiff := convertLatitudeToMeters(float64(point1.Latitude - point2.Latitude))
			pointDistance = append(pointDistance, math.Sqrt(longtitudeDiff*longtitudeDiff+latitudeDiff*latitudeDiff))
		}

	}
	return pointDistance
}
func ifNeighbours(point *mmPoint, pointDiff []float64, maxDistance float64, amoutOfPhotoes int64) {
	var i int64 = 0
	for _, value := range pointDiff {
		fmt.Println(value)
		if value <= maxDistance {
			i++
		}
	}

	if i < amoutOfPhotoes {
		point.setColor(NotValidPointR, NotValidPointG, NotValidPointB)
	}

}
