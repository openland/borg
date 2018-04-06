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

// intersectPoints = (poly, origin, alpha) ->
//   eps = 1e-9
//   origin = [origin[0] + eps*Math.cos(alpha), origin[1] + eps*Math.sin(alpha)]
//   [x0, y0] = origin
//   shiftedOrigin = [x0 + Math.cos(alpha), y0 + Math.sin(alpha)]

//   idx = 0
//   if Math.abs(shiftedOrigin[0] - x0) < eps then idx = 1
//   i = -1
//   n = poly.length
//   b = poly[n-1]
//   minSqDistLeft = Number.MAX_VALUE
//   minSqDistRight = Number.MAX_VALUE
//   closestPointLeft = null
//   closestPointRight = null
//   while ++i < n
//     a = b
//     b = poly[i]
//     p = lineIntersection origin, shiftedOrigin, a, b
//     if p? and pointInSegmentBox p, a, b
//       sqDist = squaredDist origin, p
//       if p[idx] < origin[idx]
//         if sqDist < minSqDistLeft
//           minSqDistLeft = sqDist
//           closestPointLeft = p
//       else if p[idx] > origin[idx]
//         if sqDist < minSqDistRight
//           minSqDistRight = sqDist
//           closestPointRight = p

//   return [closestPointLeft, closestPointRight]

func loadEdgedCenters(poly geometry.Polygon2D, smallSide float64, largeSide float64) []geometry.Point2D {
	res := make([]geometry.Point2D, 0)
	edges := poly.EdgesVectors()
	for _, e := range edges {
		// l := e.Length()

		// Offsetting
		offset := e.Normal().Identity().Multiply(smallSide)

		// Forward
		center := geometry.Point2D{X: e.Origin.X + offset.DX + e.DX/2, Y: e.Origin.Y + offset.DY + e.DY/2}
		res = append(res, center)
		left, right := poly.RayIntersections(center, e.DX, e.DY)
		if left != nil && right != nil {
			res = append(res, geometry.Point2D{X: (left.X + right.X - offset.DX) / 2, Y: (left.Y + right.Y - offset.DY) / 2})
		}

		// Backward
		center = geometry.Point2D{X: e.Origin.X - offset.DX + e.DX/2, Y: e.Origin.Y - offset.DY + e.DY/2}
		res = append(res, center)
		left, right = poly.RayIntersections(center, e.DX, e.DY)
		if left != nil && right != nil {
			res = append(res, geometry.Point2D{X: (left.X + right.X - offset.DX) / 2, Y: (left.Y + right.Y - offset.DY) / 2})
		}

		// fmt.Println("Edge")
		// fmt.Println(e)
		// fmt.Println(center.DebugString())
		// if left != nil && right != nil {
		// 	fmt.Println(left.DebugString())
		// 	fmt.Println(right.DebugString())
		// }

		// parallel := geometry.Vector2D{
		// 	Origin: geometry.Point2D{X: e.Origin.X + offset.DX + e.DX/2, Y: e.Origin.Y + offset.DY + e.DY/2},
		// 	DX:     e.DX,
		// 	DY:     e.DY,
		// }

		// dx := offset.X - e.Origin.X
		// dy := offset.Y - e.Origin.Y

		// shifted := e
	}
	return res
}

// func layoutEdgeAligned(poly geometry.Polygon2D, smallSide float64, largeSide float64) Layout {
// 	edges := poly.EdgesVectors()

// 	for _, e := range edges {
// 		// l := e.Length()

// 		// Offsetting
// 		offset := e.Normal().Identity().Multiply(smallSide)
// 		center := geometry.Point2D{X: e.Origin.X + offset.DX + e.DX/2, Y: e.Origin.Y + offset.DY + e.DY/2}
// 		left, right := poly.RayIntersections(center, e.DX, e.DY)
// 		if left != nil && right != nil {

// 		}
// 		fmt.Println("Edge")
// 		fmt.Println(e)
// 		fmt.Println(center)
// 		fmt.Println(left)
// 		fmt.Println(right)

// 		// parallel := geometry.Vector2D{
// 		// 	Origin: geometry.Point2D{X: e.Origin.X + offset.DX + e.DX/2, Y: e.Origin.Y + offset.DY + e.DY/2},
// 		// 	DX:     e.DX,
// 		// 	DY:     e.DY,
// 		// }

// 		// dx := offset.X - e.Origin.X
// 		// dy := offset.Y - e.Origin.Y

// 		// shifted := e
// 	}

// 	return Layout{Analyzed: true, Fits: false}
// }

func LayoutRectangle(poly geometry.Polygon2D, width float64, height float64) Layout {
	t := poly.Classify()
	smallSide := math.Min(width, height)
	largeSide := math.Max(width, height)

	// Ignore all complex polygons
	if t == geometry.TypePolygonWithHoles {
		return Layout{Analyzed: false}
	}

	// Fast-handling of rectangles
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
				Analyzed:    true,
				Fits:        true,
				HasLocation: true,
				Center:      poly.Center(),
				Angle:       mainAngle}
		}
		return Layout{Analyzed: true, Fits: false}
	}

	// Rectangle
	centers := []geometry.Point2D{poly.Center()}
	sideAngles := poly.Azimuths()
	footprint := geometry.NewSimplePolygon([]geometry.Point2D{
		{-smallSide / 2, largeSide / 2},
		{smallSide / 2, largeSide / 2},
		{smallSide / 2, -largeSide / 2},
		{-smallSide / 2, -largeSide / 2},
	})

	// Aligned to edge
	for _, c := range loadEdgedCenters(poly, smallSide, largeSide) {
		centers = append(centers, c)
	}
	for _, c := range loadEdgedCenters(poly, largeSide, smallSide) {
		centers = append(centers, c)
	}

	// fmt.Println(centers)

	// Random Centers
	// bounds := poly.Bounds()
	// for i := 0; i < 100; i++ {
	// 	x := rand.Float64()*(bounds.MaxX-bounds.MinX) + bounds.MinX
	// 	y := rand.Float64()*(bounds.MaxY-bounds.MinY) + bounds.MinY
	// 	point := geometry.Point2D{X: x, Y: y}
	// 	if poly.ContainsPoint(point) {
	// 		centers = append(centers, point)
	// 	}
	// }

	// Main Loop
	for _, center := range centers {

		//
		// Fast Search (aligned to side)
		//

		for i := 0; i < len(sideAngles); i++ {
			r := footprint.
				Rotate(sideAngles[i]).
				Shift(center)

			if poly.Contains(r) {
				return Layout{
					Analyzed: true, Fits: true, HasLocation: true,
					Center: center,
					Angle:  sideAngles[i]}
			}
		}

		//
		// Full Search
		//

		rotationIterations := 1000
		for i := 0; i < rotationIterations; i++ {
			alpha := float64(i) * math.Pi * 2 / float64(rotationIterations)
			r := footprint.
				Rotate(alpha).
				Shift(center)
			if poly.Contains(r) {
				return Layout{
					Analyzed: true, Fits: true, HasLocation: true,
					Center: center,
					Angle:  alpha}
			}
		}
	}

	return Layout{Analyzed: true}
}
