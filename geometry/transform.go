package geometry

import "math"

//
// 3D Transformations
// src: https://en.wikipedia.org/wiki/Rotation_matrix#Basic_rotations
//

func (point Point3D) PrecomputedRotateX(angleSin float64, angleCos float64) Point3D {
	rotatedX := point.X
	rotatedY := angleCos*point.Y - angleSin*point.Z
	rotatedZ := angleSin*point.Y + angleCos*point.Z
	return Point3D{rotatedX, rotatedY, rotatedZ}
}

func (point Point3D) PrecomputedRotateY(angleSin float64, angleCos float64) Point3D {
	rotatedX := angleCos*point.X + angleSin*point.Z
	rotatedY := 1 * point.Y
	rotatedZ := -angleSin*point.X + angleCos*point.Z
	return Point3D{rotatedX, rotatedY, rotatedZ}
}

func (point Point3D) PrecomputedRotateZ(angleSin float64, angleCos float64) Point3D {
	rotatedX := angleCos*point.X - angleSin*point.Y
	rotatedY := angleSin*point.X + angleCos*point.Y
	rotatedZ := point.Z
	return Point3D{rotatedX, rotatedY, rotatedZ}
}

func (point Point3D) RotateX(angle float64) Point3D {
	return point.PrecomputedRotateX(math.Sin(angle), math.Cos(angle))
}

func (point Point3D) RotateY(angle float64) Point3D {
	return point.PrecomputedRotateY(math.Sin(angle), math.Cos(angle))
}

func (point Point3D) RotateZ(angle float64) Point3D {
	return point.PrecomputedRotateZ(math.Sin(angle), math.Cos(angle))
}

func (point Point3D) Shift(delta Point3D) Point3D {
	return Point3D{point.X + delta.X, point.Y + delta.Y, point.Z + delta.Z}
}

func (point Point3D) Identity() Point3D {
	l := math.Sqrt(point.X*point.X + point.Y*point.Y + point.Z*point.Z)
	return Point3D{point.X / l, point.Y / l, point.Z / l}
}

func (point Point3D) Invert() Point3D {
	return Point3D{-point.X, -point.Y, -point.Z}
}

func (point Point3D) Miltiply(scalar float64) Point3D {
	return Point3D{point.X * scalar, point.Y * scalar, point.Z * scalar}
}

//
// 2D Rotations
//

func (point Point2D) Rotate(angle float64) Point2D {
	l := math.Sqrt(point.X*point.X + point.Y*point.Y)
	a := math.Atan2(point.X, point.Y) + angle
	return Point2D{math.Sin(a) * l, math.Cos(a) * l}
}

func (point Point2D) Shift(delta Point2D) Point2D {
	return Point2D{point.X + delta.X, point.Y + delta.Y}
}
