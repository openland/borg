package geometry

type PolygonType int32

const (
	TypeComplexPolygon   PolygonType = 1
	TypePolygonWithHoles PolygonType = 2
	TypeTriangle         PolygonType = 3
	TypeRectangle        PolygonType = 4
	TypeQuadriliteral    PolygonType = 5
	TypeConvexPolygon    PolygonType = 6
	TypeBroken           PolygonType = 7
)

func (poly Polygon2D) Classify() PolygonType {
	// Check if there are any innner circles (holes)
	if len(poly.Holes) > 0 {
		return TypePolygonWithHoles
	}

	// Here we have single polygon without holes, so simplify to one line circle
	line := poly.Polygon

	// Is broken?
	if len(line) < 4 {
		return TypeBroken
	}

	// Load angles for furher processing
	angles := poly.Angles()

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
