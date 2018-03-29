package utils

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
