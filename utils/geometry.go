package utils

import (
	"math"
)

type Point struct {
	X float64
	Y float64
}

type Bounds struct {
	MinX float64
	MinY float64
	MaxX float64
	MaxY float64
}

type Point3d struct {
	X float64
	Y float64
	Z float64
}

func vectorAngle(point1 []float64, point2 []float64, point3 []float64) float64 {
	dx1 := point2[0] - point1[0]
	dy1 := point2[1] - point1[1]
	dx2 := point3[0] - point2[0]
	dy2 := point3[1] - point2[1]
	// dot := dx1*dx2 + dy1*dy2
	cross := dx1*dy2 - dy1*dx2
	len1 := math.Sqrt((dx1*dx1 + dy1*dy1))
	len2 := math.Sqrt((dx2*dx2 + dy2*dy2))
	return grad(math.Asin(cross / (len1 * len2)))
	// return grad(math.Acos(dot / (len1 * len2)))
}

func GetAngles(line [][]float64) []float64 {
	res := make([]float64, 0)
	for i := 0; i < len(line)-1; i++ {
		i1 := i
		i2 := i + 1
		i3 := i + 2
		if i3 >= len(line) {
			i3 = 1
		}

		res = append(res, vectorAngle(line[i1], line[i2], line[i3]))
	}
	return res
}

func GetSides(line [][]float64) []float64 {
	res := make([]float64, 0)
	for i := 0; i < len(line)-1; i++ {
		l := math.Sqrt((line[i][0]-line[i+1][0])*(line[i][0]-line[i+1][0]) + (line[i][1]-line[i+1][1])*(line[i][1]-line[i+1][1]))
		res = append(res, l)
	}
	return res
}

func GetSideGlobalAngles(line [][]float64) []float64 {
	res := make([]float64, 0)
	for i := 0; i < len(line)-1; i++ {
		res = append(res, GlobalAngle(line[i], line[i+1]))
	}
	return res
}

func FindBounds(coords [][][][]float64) Bounds {
	maxX := -math.MaxFloat64
	minX := math.MaxFloat64
	maxY := -math.MaxFloat64
	minY := math.MaxFloat64
	for _, poly := range coords {
		for _, circle := range poly {
			for _, point := range circle {
				if point[0] > maxX {
					maxX = point[0]
				}
				if point[0] < minX {
					minX = point[0]
				}
				if point[1] > maxY {
					maxY = point[1]
				}
				if point[1] < minY {
					minY = point[1]
				}
			}
		}
	}
	return Bounds{MinX: minX, MaxX: maxX, MinY: minY, MaxY: maxY}
}

func FindCenter(coords [][][][]float64) Point {
	bounds := FindBounds(coords)
	return Point{X: (bounds.MaxX + bounds.MinX) / 2, Y: (bounds.MaxY + bounds.MinY) / 2}
}

func IsPointsInside(points [][]float64, polygon [][]float64) bool {
	for _, p := range points {
		if !IsPointInside(p, polygon) {
			return false
		}
	}
	return true
}

func IsPointInside(point []float64, polygon [][]float64) bool {

	// http://www.ecse.rpi.edu/Homepages/wrf/Research/Short_Notes/pnpoly.html
	// https://stackoverflow.com/questions/217578/how-can-i-determine-whether-a-2d-point-is-within-a-polygon/17490923#17490923
	// https://github.com/JamesMilnerUK/pip-go/blob/master/pip.go

	isInside := false
	j := 0
	for i := 1; i < len(polygon); i++ {

		// Do not touch!
		if ((polygon[i][1] > point[1]) != (polygon[j][1] > point[1])) &&
			(point[0] < (polygon[j][0]-polygon[i][0])*(point[1]-polygon[i][1])/(polygon[j][1]-polygon[i][1])+polygon[i][0]) {
			isInside = !isInside
		}

		j = i
	}

	return isInside
}

func GlobalAngle(a, b []float64) float64 {
	return math.Atan2(b[0]-a[0], b[1]-a[1])
}

func Rotate2D(points [][]float64, angle float64) [][]float64 {
	res := make([][]float64, 0)
	for _, p := range points {
		l := math.Sqrt(p[0]*p[0] + p[1]*p[1])
		a := math.Atan2(p[0], p[1]) + angle
		res = append(res, []float64{math.Sin(a) * l, math.Cos(a) * l})
	}
	return res
}

func Shift2D(points [][]float64, shift []float64) [][]float64 {
	res := make([][]float64, 0)
	for _, p := range points {
		res = append(res, []float64{p[0] + shift[0], p[1] + shift[1]})
	}
	return res
}
