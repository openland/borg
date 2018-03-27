package utils

import (
	"math"
)

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
