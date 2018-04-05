package geometry

import "math"

func (src PointGeo) ToMercator() Point2D {
	x := src.Longitude * 20037508.34 / 180.0
	y := math.Log(math.Tan((90+src.Latitude)*math.Pi/360.0)) / (math.Pi / 180)
	y = y * 20037508.34 / 180.0
	return Point2D{X: x, Y: y}
}
