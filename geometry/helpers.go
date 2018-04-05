package geometry

import "math"

const eps = 1e-9

func rad(src float64) float64 {
	return src * math.Pi / 180
}

func grad(src float64) float64 {
	return src * 180 / math.Pi
}
