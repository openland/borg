package geometry

import (
	"fmt"
	"math"
)

//
// Points
//

// Point2D is a point in 2D coordinates
type Point2D struct {
	X float64
	Y float64
}

// Point3D is a point in 3D coordinates
type Point3D struct {
	X float64
	Y float64
	Z float64
}

// PointGeo is a point in Geo coordinates
type PointGeo struct {
	Longitude float64
	Latitude  float64
}

func (point Point2D) Distance(to Point2D) float64 {
	return math.Sqrt((point.X-to.X)*(point.X-to.X) + (point.Y-to.Y)*(point.Y-to.Y))
}

func (point Point2D) DistanceSq(to Point2D) float64 {
	return (point.X-to.X)*(point.X-to.X) + (point.Y-to.Y)*(point.Y-to.Y)
}

func (point Point2D) Azimuth(to Point2D) float64 {
	return math.Atan2(to.X-point.X, to.Y-point.Y)
}

func (point Point2D) DebugString() string {
	return fmt.Sprintf("(%.6f,%.6f)", point.X, point.Y)
}
