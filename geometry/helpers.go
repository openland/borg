package geometry

import "math"

func rad(src float64) float64 {
	return src * math.Pi / 180
}

func grad(src float64) float64 {
	return src * 180 / math.Pi
}
