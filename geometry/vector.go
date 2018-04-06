package geometry

import (
	"math"
)

type Vector2D struct {
	Origin Point2D
	DX     float64
	DY     float64
}

func (vector Vector2D) LengthSq() float64 {
	return (vector.DX*vector.DX + vector.DY*vector.DY)
}

func (vector Vector2D) Length() float64 {
	return math.Sqrt(vector.LengthSq())
}

func (vector Vector2D) Normal() Vector2D {
	return Vector2D{Origin: vector.Origin, DX: vector.DY, DY: -vector.DX}
}

func (vector Vector2D) Identity() Vector2D {
	l := vector.Length()
	return Vector2D{Origin: vector.Origin, DX: vector.DX / l, DY: vector.DY / l}
}

func (vector Vector2D) Multiply(v float64) Vector2D {
	return Vector2D{Origin: vector.Origin, DX: vector.DX * v, DY: vector.DY * v}
}

func (vector Vector2D) DebugString() string {
	return vector.Origin.DebugString() + "-" + Point2D{X: vector.Origin.X + vector.DX, Y: vector.Origin.Y + vector.DY}.DebugString()
}
