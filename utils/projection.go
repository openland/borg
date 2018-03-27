package utils

import (
	"fmt"
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

func ProjectMercator(latitude float64, longitude float64) Point {
	x := longitude * 20037508.34 / 180.0
	y := math.Log(math.Tan((90+latitude)*math.Pi/360.0)) / (math.Pi / 180)
	y = y * 20037508.34 / 180.0
	return Point{X: x, Y: y}
}

func ProjectToPlane(coords [][][][]float64) [][][][]float64 {
	res := make([][][][]float64, 0)
	for _, poly := range coords {
		cpoly := make([][][]float64, 0)
		for _, circle := range poly {
			ccircle := make([][]float64, 0)
			for _, point := range circle {
				cpoint := ProjectMercator(point[1], point[0])
				ccircle = append(ccircle, []float64{cpoint.X, cpoint.Y})
			}
			cpoly = append(cpoly, ccircle)
		}
		res = append(res, cpoly)
	}
	bounds := FindBounds(res)
	fmt.Println(bounds)
	for _, poly := range res {
		for _, circle := range poly {
			for _, point := range circle {
				point[0] = point[0] - bounds.MinX
				point[1] = point[1] - bounds.MinY
			}
		}
	}
	return res
}
