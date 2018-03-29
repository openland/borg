package utils

import (
	"math"
)

type Projection struct {
	center  Point
	centerX float64
	centerY float64
	centerZ float64

	cosLon  float64
	sinLon  float64
	cosNLon float64
	sinNLon float64

	cosLat  float64
	sinLat  float64
	cosNLat float64
	sinNLat float64
}

func NewProjection(plane [][][][]float64) *Projection {
	center := FindCenter(plane)
	centerProjected := ProjectToCartesian([]float64{center.X, center.Y})
	centerX := centerProjected[0]
	centerY := centerProjected[1]
	centerZ := centerProjected[2]

	// Avoid multiple calculations of the same values
	cosLon := math.Cos(rad(-center.X))
	sinLon := math.Sin(rad(-center.X))
	cosLat := math.Cos(rad(center.Y))
	sinLat := math.Sin(rad(center.Y))

	cosNLon := math.Cos(rad(center.X))
	sinNLon := math.Sin(rad(center.X))
	cosNLat := math.Cos(rad(-center.Y))
	sinNLat := math.Sin(rad(-center.Y))

	return &Projection{
		center:  center,
		centerX: centerX,
		centerY: centerY,
		centerZ: centerZ,

		cosLon: cosLon,
		sinLon: sinLon,
		cosLat: cosLat,
		sinLat: sinLat,

		cosNLon: cosNLon,
		sinNLon: sinNLon,
		cosNLat: cosNLat,
		sinNLat: sinNLat,
	}
}

func (proj *Projection) ProjectPoint(point []float64) []float64 {

	//
	// Projecting of point to a projection plane
	//

	// Plane equation:
	// X * centerX + Y * centerY + Z * centerZ = centerX * centerX + centerY * centerY + centerZ * centerZ
	// Line equations:
	// x = t * pointX
	// y = t * pointY
	// z = t * pointZ
	// Intersection
	// t = (centerX * centerX + centerY * centerY + centerZ * centerZ) / (pointX * centerX + pointY * centerY + pointZ * pointZ)

	pointProjected := ProjectToCartesian(point)
	pointX := pointProjected[0]
	pointY := pointProjected[1]
	pointZ := pointProjected[2]
	t := (proj.centerX*proj.centerX + proj.centerY*proj.centerY + proj.centerZ*proj.centerZ) / (pointX*proj.centerX + pointY*proj.centerY + pointZ*proj.centerZ)
	projectedX := t * pointX
	projectedY := t * pointY
	projectedZ := t * pointZ

	//
	// Move center to center of coordinates
	//
	projectedX = projectedX - proj.centerX
	projectedY = projectedY - proj.centerY
	projectedZ = projectedZ - proj.centerZ

	//
	// Rotate Plane to match XY coordinates
	//
	// Rotate around vertical axis (Z)
	rotated := RotateZ([]float64{projectedX, projectedY, projectedZ}, proj.sinLon, proj.cosLon)
	// Rotsate around horizontal axis (Y)
	rotated = RotateY(rotated, proj.sinLat, proj.cosLat)

	// Ignoring X coordinate as is is aligned with normal to our plane (and should be zero!)
	return []float64{rotated[1], rotated[2]}
}

func (proj *Projection) UnprojectPoint(point []float64) []float64 {

	//
	// Rotate coorinates back
	//

	rotated := []float64{0, point[0], point[1]}

	// Rotsate around horizontal axis (Y)
	rotated = RotateY(rotated, proj.sinNLat, proj.cosNLat)
	// Rotate around vertical axis (Z)
	rotated = RotateZ([]float64{0, point[0], point[1]}, proj.sinNLon, proj.cosNLon)

	//
	// Shift center back
	//

	projectedX := rotated[0] + proj.centerX
	projectedY := rotated[1] + proj.centerY
	projectedZ := rotated[2] + proj.centerZ

	//
	// Normalize
	//

	l := math.Sqrt(projectedX*projectedX + projectedY*projectedY + projectedZ*projectedZ)
	projectedX = WORLD_RADIUS * (projectedX / l)
	projectedY = WORLD_RADIUS * (projectedY / l)
	projectedZ = WORLD_RADIUS * (projectedZ / l)

	//
	// Unproject from cartesian to latitude/longitude
	//

	return UnprojectFromCartesian([]float64{projectedX, projectedY, projectedZ})
}

func (proj *Projection) ProjectMultiPolygon(coords [][][][]float64) [][][][]float64 {
	res := make([][][][]float64, 0)
	for _, poly := range coords {
		cpoly := make([][][]float64, 0)
		for _, circle := range poly {
			ccircle := make([][]float64, 0)
			for _, point := range circle {
				ccircle = append(ccircle, proj.ProjectPoint(point))
			}
			cpoly = append(cpoly, ccircle)
		}
		res = append(res, cpoly)
	}
	return res
}

func (proj *Projection) UnprojectMultiPolygon(coords [][][][]float64) [][][][]float64 {
	res := make([][][][]float64, 0)
	for _, poly := range coords {
		cpoly := make([][][]float64, 0)
		for _, circle := range poly {
			ccircle := make([][]float64, 0)
			for _, point := range circle {
				ccircle = append(ccircle, proj.UnprojectPoint(point))
			}
			cpoly = append(cpoly, ccircle)
		}
		res = append(res, cpoly)
	}
	return res
}

//
// Mercator Projection
//

func ProjectMercator(latitude float64, longitude float64) Point {
	x := longitude * 20037508.34 / 180.0
	y := math.Log(math.Tan((90+latitude)*math.Pi/360.0)) / (math.Pi / 180)
	y = y * 20037508.34 / 180.0
	return Point{X: x, Y: y}
}

//
// Spherical to Cartesian Projections
//

func ProjectToCartesian(src []float64) []float64 {
	centerX := WORLD_RADIUS * math.Cos(rad(src[1])) * math.Cos(rad(src[0]))
	centerY := WORLD_RADIUS * math.Cos(rad(src[1])) * math.Sin(rad(src[0]))
	centerZ := WORLD_RADIUS * math.Sin(rad(src[1]))
	return []float64{centerX, centerY, centerZ}
}

func UnprojectFromCartesian(src []float64) []float64 {
	latitude := math.Asin(src[2]/WORLD_RADIUS) * 180 / math.Pi
	var longitude float64
	if src[0] > 0 {
		longitude = math.Atan(src[1]/src[0]) * 180 / math.Pi
	} else if src[1] > 0 {
		longitude = math.Atan(src[1]/src[0])*180/math.Pi + 180
	} else {
		longitude = math.Atan(src[1]/src[0])*180/math.Pi - 180
	}
	return []float64{longitude, latitude}
}
