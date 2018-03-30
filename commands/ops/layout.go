package ops

import (
	"math"

	"github.com/statecrafthq/borg/utils"
)

type Layout struct {
	Analyzed    bool
	Fits        bool
	HasLocation bool
	Center      []float64
	Angle       float64
}

func LayoutRectangle(polys [][][][]float64, width float64, height float64) Layout {
	t := ClassifyParcelGeometry(polys)
	smallSide := math.Min(width, height)
	largeSide := math.Max(width, height)
	// Ignore all complex polygons
	if t == TypeMultipolygon || t == TypePolygonWithHoles || t == TypeComplexPolygon {
		return Layout{Analyzed: false}
	}
	poly := polys[0][0]

	// Rectangle
	if t == TypeRectangle {
		sides := utils.GetSides(poly)
		center := utils.FindCenter(polys)
		side1 := (sides[0] + sides[2]) / 2
		side2 := (sides[1] + sides[3]) / 2
		// Second side is inverted to make them aligned
		angle1 := (utils.GlobalAngle(poly[0], poly[1]) + utils.GlobalAngle(poly[3], poly[2])) / 2
		angle2 := (utils.GlobalAngle(poly[1], poly[2]) + utils.GlobalAngle(poly[4], poly[3])) / 2
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
				Center: []float64{center.X, center.Y},
				Angle:  mainAngle}
		}
		return Layout{Analyzed: true, Fits: false}
	}

	//
	// Convex Polygon: Pick center and try aligned with any side of polygon
	//
	sideAngles := utils.GetSideGlobalAngles(poly)
	center := utils.FindCenter(polys)
	baseRect := [][]float64{
		{-largeSide / 2, smallSide / 2},
		{largeSide / 2, smallSide / 2},
		{largeSide / 2, -smallSide / 2},
		{-largeSide / 2, -smallSide / 2},
		{-largeSide / 2, smallSide / 2}}

	for i := 0; i < len(sideAngles); i++ {
		rect := utils.Rotate2D(baseRect, sideAngles[i])
		rect = utils.Shift2D(rect, []float64{center.X, center.Y})
		if utils.IsPointsInside(rect, poly) {
			return Layout{
				Analyzed: true, Fits: true, HasLocation: true,
				Center: []float64{center.X, center.Y},
				Angle:  sideAngles[i]}
		}
	}

	return Layout{Analyzed: true}

}
