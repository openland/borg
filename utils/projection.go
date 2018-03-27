package utils

import (
	"math"
)

type Bounds struct {
	MinX float64
	MinY float64
	MaxX float64
	MaxY float64
}

type Point struct {
	X float64
	Y float64
}

type Point3d struct {
	X float64
	Y float64
	Z float64
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

func ProjectMercator(latitude float64, longitude float64) Point {
	x := longitude * 20037508.34 / 180.0
	y := math.Log(math.Tan((90+latitude)*math.Pi/360.0)) / (math.Pi / 180)
	y = y * 20037508.34 / 180.0
	return Point{X: x, Y: y}
}

func ProjectToCartesian(src []float64) []float64 {
	centerX := 6378137 * math.Cos(rad(src[1])) * math.Cos(rad(src[0]))
	centerY := 6378137 * math.Cos(rad(src[1])) * math.Sin(rad(src[0]))
	centerZ := 6378137 * math.Sin(rad(src[1]))
	return []float64{centerX, centerY, centerZ}
}

func RotateX(src []float64, angleSin float64, angleCos float64) []float64 {
	// https://en.wikipedia.org/wiki/Rotation_matrix#Basic_rotations
	rotatedX := 1*src[0] + 0*src[1] + 0*src[2]
	rotatedY := 0*src[0] + angleCos*src[1] - angleSin*src[2]
	rotatedZ := 0*src[0] + angleSin*src[1] + angleCos*src[2]
	return []float64{rotatedX, rotatedY, rotatedZ}
}

func RotateY(src []float64, angleSin float64, angleCos float64) []float64 {
	// https://en.wikipedia.org/wiki/Rotation_matrix#Basic_rotations
	rotatedX := angleCos*src[0] + 0*src[1] + angleSin*src[2]
	rotatedY := 0*src[0] + 1*src[1] + 0*src[2]
	rotatedZ := -angleSin*src[0] + 0*src[1] + angleCos*src[2]
	return []float64{rotatedX, rotatedY, rotatedZ}
}

func RotateZ(src []float64, angleSin float64, angleCos float64) []float64 {
	// https://en.wikipedia.org/wiki/Rotation_matrix#Basic_rotations
	rotatedX := angleCos*src[0] - angleSin*src[1] + 0*src[2]
	rotatedY := angleSin*src[0] + angleCos*src[1] + 0*src[2]
	rotatedZ := 0*src[0] + 0*src[1] + 1*src[2]
	return []float64{rotatedX, rotatedY, rotatedZ}
}

func ProjectToPlane(coords [][][][]float64) [][][][]float64 {
	res := make([][][][]float64, 0)
	center := FindCenter(coords)
	centerProjected := ProjectToCartesian([]float64{center.X, center.Y})
	centerX := centerProjected[0]
	centerY := centerProjected[1]
	centerZ := centerProjected[2]

	// Avoid multiple calculations of the same values
	cosLon := math.Cos(rad(-center.X))
	sinLon := math.Sin(rad(-center.X))
	cosLat := math.Cos(rad(center.Y))
	sinLat := math.Sin(rad(center.Y))

	// Rotation of Center (test code)
	// rotatedCenter := RotateZ(centerProjected, sinLon, cosLon)
	// fmt.Printf("Center(1): %.6f, %.6f, %.6f\n", rotatedCenter[0], rotatedCenter[1], rotatedCenter[2])
	// rotatedCenter = RotateY(rotatedCenter, sinLat, cosLat)
	// fmt.Printf("Center(2): %.6f, %.6f, %.6f\n", rotatedCenter[0], rotatedCenter[1], rotatedCenter[2])

	for _, poly := range coords {
		cpoly := make([][][]float64, 0)
		for _, circle := range poly {
			ccircle := make([][]float64, 0)
			for _, point := range circle {

				//
				// Projecting of point to a plane and shift co center of coordiantes
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
				t := (centerX*centerX + centerY*centerY + centerZ*centerZ) / (pointX*centerX + pointY*centerY + pointZ*centerZ)
				projectedX := t*pointX - centerX
				projectedY := t*pointY - centerY
				projectedZ := t*pointZ - centerZ

				//
				// Rotate Plane to match XY coordinates
				//

				// Rotate around vertical axis (Z)
				rotated := RotateZ([]float64{projectedX, projectedY, projectedZ}, sinLon, cosLon)
				// Rotsate around horizontal axis (Y)
				rotated = RotateY(rotated, sinLat, cosLat)

				// Ignoring X coordinate as is is aligned with normal to our plane
				ccircle = append(ccircle, []float64{rotated[1], rotated[2]})
			}
			cpoly = append(cpoly, ccircle)
		}
		res = append(res, cpoly)
	}
	return res
}

// func ProjectToPlane(coords [][][][]float64) [][][][]float64 {
// 	res := make([][][][]float64, 0)
// 	for _, poly := range coords {
// 		cpoly := make([][][]float64, 0)
// 		for _, circle := range poly {
// 			ccircle := make([][]float64, 0)
// 			for _, point := range circle {
// 				cpoint := ProjectMercator(point[1], point[0])
// 				ccircle = append(ccircle, []float64{cpoint.X, cpoint.Y})
// 			}
// 			cpoly = append(cpoly, ccircle)
// 		}
// 		res = append(res, cpoly)
// 	}
// 	bounds := FindBounds(res)
// 	fmt.Println(bounds)
// 	for _, poly := range res {
// 		for _, circle := range poly {
// 			for _, point := range circle {
// 				point[0] = point[0] - bounds.MinX
// 				point[1] = point[1] - bounds.MinY
// 			}
// 		}
// 	}
// 	return res
// }
