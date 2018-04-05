package geometry

import (
	"fmt"
	"math"
)

// LineString is a sequence of 2d points
type LineString2D = []Point2D

// Polygon in 2D space
type Polygon2D struct {
	Polygon LineString2D
	Holes   []LineString2D
}

// Multipolygon in 2D space
type Multipolygon2D struct {
	Polygons []Polygon2D
}

func NewProjectedSimplePolygon(points []PointGeo) Polygon2D {
	p := NewSimpleGeoPolygon(points)
	return p.Project(NewProjection(p.Center()))
}

func NewSimplePolygon(points []Point2D) Polygon2D {
	return Polygon2D{Polygon: points, Holes: [][]Point2D{}}
}

func (poly Polygon2D) Edges() []float64 {
	res := make([]float64, 0)
	for i := 0; i < len(poly.Polygon); i++ {
		s1 := poly.Polygon[i]
		s2 := poly.Polygon[(i+1)%len(poly.Polygon)]
		res = append(res, s1.Distance(s2))
	}
	return res
}

func vectorAngle(point1 Point2D, point2 Point2D, point3 Point2D) float64 {
	dx1 := point2.X - point1.X
	dy1 := point2.Y - point1.Y
	dx2 := point3.X - point2.X
	dy2 := point3.Y - point2.Y
	// dot := dx1*dx2 + dy1*dy2
	cross := dx1*dy2 - dy1*dx2
	len1 := math.Sqrt((dx1*dx1 + dy1*dy1))
	len2 := math.Sqrt((dx2*dx2 + dy2*dy2))
	return grad(math.Asin(cross / (len1 * len2)))
	// return grad(math.Acos(dot / (len1 * len2)))
}

func (poly Polygon2D) Angles() []float64 {
	res := make([]float64, 0)
	for i := 0; i < len(poly.Polygon); i++ {
		i1 := i
		i2 := (i + 1) % len(poly.Polygon)
		i3 := (i + 2) % len(poly.Polygon)
		res = append(res, vectorAngle(poly.Polygon[i1], poly.Polygon[i2], poly.Polygon[i3]))
	}
	return res
}

func (poly Polygon2D) Azimuths() []float64 {
	res := make([]float64, 0)
	for i := 0; i < len(poly.Polygon); i++ {
		i1 := i
		i2 := (i + 1) % len(poly.Polygon)
		res = append(res, poly.Polygon[i1].Azimuth(poly.Polygon[i2]))
	}
	return res
}

func (polygon Polygon2D) Center() Point2D {
	count := 0
	xS := 0.0
	yS := 0.0
	for i := 0; i < len(polygon.Polygon); i++ {
		count++
		xS += polygon.Polygon[i].X
		yS += polygon.Polygon[i].Y
	}
	return Point2D{X: xS / float64(count), Y: yS / float64(count)}
}

func (polygon Polygon2D) Rotate(alpha float64) Polygon2D {
	main := make([]Point2D, 0)
	for _, p := range polygon.Polygon {
		main = append(main, p.Rotate(alpha))
	}
	holes := make([][]Point2D, 0)
	for _, h := range polygon.Holes {
		p2 := make([]Point2D, 0)
		for _, p3 := range h {
			p2 = append(p2, p3.Rotate(alpha))
		}
		holes = append(holes, p2)
	}
	return Polygon2D{Polygon: main, Holes: holes}
}

func (polygon Polygon2D) Shift(delta Point2D) Polygon2D {
	main := make([]Point2D, 0)
	for _, p := range polygon.Polygon {
		main = append(main, p.Shift(delta))
	}
	holes := make([][]Point2D, 0)
	for _, h := range polygon.Holes {
		p2 := make([]Point2D, 0)
		for _, p3 := range h {
			p2 = append(p2, p3.Shift(delta))
		}
		holes = append(holes, p2)
	}
	return Polygon2D{Polygon: main, Holes: holes}
}

func (polygon Polygon2D) Contains(dst Polygon2D) bool {
	for _, p := range dst.Polygon {
		if !polygon.ContainsPoint(p) {
			return false
		}
	}
	return true
}

func (polygon Polygon2D) ContainsAllPoints(points []Point2D) bool {
	for _, p := range points {
		if !polygon.ContainsPoint(p) {
			return false
		}
	}
	return true
}

func containsPoint(point Point2D, lineString LineString2D) bool {
	// http://www.ecse.rpi.edu/Homepages/wrf/Research/Short_Notes/pnpoly.html
	// https://stackoverflow.com/questions/217578/how-can-i-determine-whether-a-2d-point-is-within-a-polygon/17490923#17490923
	// https://github.com/JamesMilnerUK/pip-go/blob/master/pip.go

	isInside := false
	j := 0
	for i := 0; i < len(lineString); i++ {

		// Do not touch!
		if ((lineString[i].Y > point.Y) != (lineString[j].Y > point.Y)) &&
			(point.X < (lineString[j].X-lineString[i].X)*(point.Y-lineString[i].Y)/(lineString[j].Y-lineString[i].Y)+lineString[i].X) {
			isInside = !isInside
		}

		j = i
	}

	return isInside
}

func (polygon Polygon2D) ContainsPoint(point Point2D) bool {

	// If not in polygon return
	if !containsPoint(point, polygon.Polygon) {
		return false
	}

	// Check if point is actually in hole
	for _, h := range polygon.Holes {
		if containsPoint(point, h) {
			return false
		}
	}

	return true
}

func (poly Polygon2D) DebugString() string {
	res := ""
	isFirst := true
	for _, point := range poly.Polygon {
		// Ignore last
		if isFirst {
			isFirst = false
		} else {
			res = res + ","
		}
		res = res + fmt.Sprintf("(%.6f,%.6f)", point.X, point.Y)
	}
	return res
}
