package ops

import (
	"math"

	"github.com/statecrafthq/borg/geometry"
)

type Layout struct {
	Analyzed    bool
	Fits        bool
	HasLocation bool
	Center      geometry.Point2D
	Angle       float64
}

func LayoutRectangle(poly geometry.Polygon2D, width float64, height float64) Layout {
	t := poly.Classify()
	smallSide := math.Min(width, height)
	largeSide := math.Max(width, height)
	center := poly.Center()

	// Ignore all complex polygons
	if t == geometry.TypePolygonWithHoles || t == geometry.TypeComplexPolygon {
		return Layout{Analyzed: false}
	}

	// Rectangle
	if t == geometry.TypeRectangle {
		sides := poly.Edges()
		side1 := (sides[0] + sides[2]) / 2
		side2 := (sides[1] + sides[3]) / 2
		// Second side is inverted to make them aligned
		angle1 := (poly.Polygon[0].Azimuth(poly.Polygon[1]) + poly.Polygon[3].Azimuth(poly.Polygon[2])) / 2
		angle2 := (poly.Polygon[1].Azimuth(poly.Polygon[2]) + poly.Polygon[0].Azimuth(poly.Polygon[3])) / 2
		var mainAngle float64
		if side1 > side2 {
			mainAngle = angle1
		} else {
			mainAngle = angle2
		}
		small := math.Min((sides[0]+sides[2])/2, (sides[1]+sides[3])/2)
		large := math.Max((sides[0]+sides[2])/2, (sides[1]+sides[3])/2)

		if small > smallSide && large > largeSide {
			return Layout{
				Analyzed: true, Fits: true, HasLocation: true,
				Center: center,
				Angle:  mainAngle}
		}
		return Layout{Analyzed: true, Fits: false}
	}

	//
	// Convex Polygon: Pick center and try aligned with any side of polygon
	//
	sideAngles := poly.Azimuths()

	rect := geometry.NewSimplePolygon(
		[]geometry.Point2D{
			{-smallSide / 2, largeSide / 2},
			{smallSide / 2, largeSide / 2},
			{smallSide / 2, -largeSide / 2},
			{-smallSide / 2, -largeSide / 2},
		})

	for i := 0; i < len(sideAngles); i++ {
		r := rect.
			Rotate(sideAngles[i]).
			Shift(center)

		if poly.Contains(r) {
			return Layout{
				Analyzed: true, Fits: true, HasLocation: true,
				Center: center,
				Angle:  sideAngles[i]}
		}
	}

	return Layout{Analyzed: true}
}
