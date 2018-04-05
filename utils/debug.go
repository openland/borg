package utils

// import (
// 	"fmt"

// 	"github.com/statecrafthq/borg/geometry"
// )

// func DebugPointLine(circles []geometry.Point2D) string {
// 	res := ""
// 	isFirst := true
// 	for i, point := range circles {
// 		// Ignore last
// 		if i == len(circles)-1 {
// 			continue
// 		}
// 		if isFirst {
// 			isFirst = false
// 		} else {
// 			res = res + ","
// 		}
// 		res = res + fmt.Sprintf("(%.6f,%.6f)", point.X, point.Y)
// 	}
// 	return res
// }

// func DebugPointLine3(circles []geometry.Point3D) string {
// 	res := ""
// 	isFirst := true
// 	for i, point := range circles {
// 		// Ignore last
// 		if i == len(circles)-1 {
// 			continue
// 		}
// 		if isFirst {
// 			isFirst = false
// 		} else {
// 			res = res + ","
// 		}
// 		res = res + fmt.Sprintf("(%.6f,%.6f,%.6f)", point.X, point.Y, point.Z)
// 	}
// 	return res
// }

// func DebugPolygon(circles [][]geometry.Point2D) string {
// 	if len(circles) <= 1 {
// 		return DebugPointLine(circles[0])
// 	} else {
// 		res := "("
// 		isFirst := true
// 		for _, cicle := range circles {
// 			if isFirst {
// 				isFirst = false
// 			} else {
// 				res = res + ","
// 			}
// 			res = res + DebugPointLine(cicle)
// 		}
// 		return res
// 	}
// }

// func DebugMultiPolygon(polys [][][]geometry.Point3D) string {
// 	if len(polys) <= 1 {
// 		return DebugPolygon(polys[0])
// 	} else {
// 		res := "("
// 		isFirst := true
// 		for _, poly := range polys {
// 			if isFirst {
// 				isFirst = false
// 			} else {
// 				res = res + ","
// 			}
// 			res = res + DebugPolygon(poly)
// 		}
// 		return res
// 	}
// }
