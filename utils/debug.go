package utils

import (
	"fmt"
)

func DebugPointLine(circles [][]float64) string {
	res := ""
	isFirst := true
	for i, point := range circles {
		// Ignore last
		if i == len(circles)-1 {
			continue
		}
		if isFirst {
			isFirst = false
		} else {
			res = res + ","
		}
		res = res + fmt.Sprintf("(%.6f,%.6f)", point[0], point[1])
	}
	return res
}

func DebugPolygon(circles [][][]float64) string {
	if len(circles) <= 1 {
		return DebugPointLine(circles[0])
	} else {
		res := "("
		isFirst := true
		for _, cicle := range circles {
			if isFirst {
				isFirst = false
			} else {
				res = res + ","
			}
			res = res + DebugPointLine(cicle)
		}
		return res
	}
}

func DebugMultiPolygon(polys [][][][]float64) string {
	if len(polys) <= 1 {
		return DebugPolygon(polys[0])
	} else {
		res := "("
		isFirst := true
		for _, poly := range polys {
			if isFirst {
				isFirst = false
			} else {
				res = res + ","
			}
			res = res + DebugPolygon(poly)
		}
		return res
	}
}
