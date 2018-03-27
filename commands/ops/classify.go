package ops

import "github.com/statecrafthq/borg/utils"

type ParcelType int32

const (
	TypeMultipolygon     ParcelType = 0
	TypeComplexPolygon   ParcelType = 1
	TypePolygonWithHoles ParcelType = 2
	TypeTriangle         ParcelType = 3
	TypeRectangle        ParcelType = 4
	TypeQuadriliteral    ParcelType = 5
	TypeConvexPolygon    ParcelType = 6
	TypeBroken           ParcelType = 7
)

func ClassifyParcelGeometry(projected [][][][]float64) ParcelType {

	// If more than one polygon mark is as multipolygon and do not analyze further
	if len(projected) > 1 {
		return TypeMultipolygon
	}

	// Check if there are any innner circles (holes)
	if len(projected[0]) > 1 {
		return TypePolygonWithHoles
	}

	// Here we have single polygon without holes, so simplify to one line circle
	line := projected[0][0]

	// Is broken?
	if len(line) < 4 {
		return TypeBroken
	}

	// Load angles for furher processing
	angles := utils.GetAngles(line)

	// Check if polygon is complex (eg NOT convex)
	for _, angle := range angles {
		if angle < 0 {
			return TypeComplexPolygon
		}
	}

	// Check if our polygon is triangle
	if len(line) == 4 {
		return TypeTriangle
	}

	// Handle Quadriliteral case
	if len(line) == 5 {

		// Check if polygon is rectangle
		isRect := true
		for _, angle := range angles {
			if angle > 91 || angle < 89 {
				isRect = false
				break
			}
		}
		if isRect {
			return TypeRectangle
		}

		// Just generic Quadriliteral
		return TypeQuadriliteral
	}

	// Generic convex polygon
	return TypeConvexPolygon
}
