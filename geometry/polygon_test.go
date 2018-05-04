package geometry

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAngles(t *testing.T) {
	testPoly := NewSimplePolygon([]Point2D{
		{-8.690199761239422, 1.391498730615154},
		{5.231000003567562, -7.068785818689552},
		{8.690202895388282, -1.2801690473393077},
		{-5.230990022569944, 7.068789511576092},
	})
	res := testPoly.Angles()

	// Number of Angles should be equals to number of unique points
	assert.Equal(t, 4, len(res))

	// In this case all angles are about 90 degree
	for _, angl := range res {
		assert.InEpsilon(t, 90, angl, 0.1)
	}

	// Random complex polygon from NYC
	irregular := NewProjectedSimplePolygon([]PointGeo{
		{-74.00591, 40.719953},
		{-74.005984, 40.719981},
		{-74.005982, 40.719993},
		{-74.006059, 40.720001},
		{-74.006066, 40.719962},
		{-74.005942, 40.719912},
	})
	res = irregular.Angles()

	// One angle should be negative as it is a complex polygon
	assert.InEpsilon(t, -70, res[0], 0.1)

	// Almost rectangle
	poly2 := NewProjectedSimplePolygon([]PointGeo{
		{-73.95109, 40.703115},
		{-73.951331, 40.70333},
		{-73.951413, 40.703276},
		{-73.951173, 40.703062},
		{-73.951159, 40.70307},
	})
	res = poly2.Angles()

	assert.InEpsilon(t, 89, res[0], 1)
	assert.InEpsilon(t, 89, res[1], 1)
	assert.InEpsilon(t, 86, res[2], 1)
	assert.InEpsilon(t, 3, res[3], 1)
	assert.InEpsilon(t, 89, res[4], 1)
}

func TestEdges(t *testing.T) {
	testPoly := NewSimplePolygon([]Point2D{
		{-8.690199761239422, 1.391498730615154},
		{5.231000003567562, -7.068785818689552},
		{8.690202895388282, -1.2801690473393077},
		{-5.230990022569944, 7.068789511576092}})
	res := testPoly.Edges()
	assert.InEpsilon(t, 16.2903, res[0].Length(), 0.0001)
	assert.InEpsilon(t, 6.7434, res[1].Length(), 0.0001)
	assert.InEpsilon(t, 16.2328, res[2].Length(), 0.0001)
	assert.InEpsilon(t, 6.6481, res[3].Length(), 0.0001)
}

func TestGlobalAngle(t *testing.T) {

	// 45 degres angle
	angle := Point2D{0, 0}.Azimuth(Point2D{1, 1})
	assert.InEpsilon(t, math.Pi/4, angle, 0.000001)
	angle = Point2D{1, 1}.Azimuth(Point2D{2, 2})
	assert.InEpsilon(t, math.Pi/4, angle, 0.000001)

	// -45 degres angle
	angle = Point2D{0, 0}.Azimuth(Point2D{-1, 1})
	assert.InEpsilon(t, -math.Pi/4, angle, 0.000001)
	angle = Point2D{1, 1}.Azimuth(Point2D{0, 2})
	assert.InEpsilon(t, -math.Pi/4, angle, 0.000001)
}

func TestPointTest(t *testing.T) {
	// WP notes a ray passing through a "side" vertex is an interesting test case.
	// The test case selected for ExampleXY shows the function working properly
	// in this case.
	res := NewSimplePolygon([]Point2D{{0, 0}, {0, 2}, {2, 1}}).ContainsPoint(Point2D{1, 1})
	assert.Equal(t, true, res)

	// Some random polygon
	poly := NewSimplePolygon([]Point2D{{17.926384, -31.113776}, {35.388660, -3.729118}, {-18.263668, 31.113820}, {-35.388622, 4.174565}, {-9.068620, -13.525313}})
	res = poly.ContainsPoint(Point2D{X: 0, Y: 0})
	assert.Equal(t, true, res)

	//[-15.18718918268496 -1.558457309876673] [-2.1937057966535596 -3.896181854342805] [10.884159114175638 -6.233883493017709]]]]
	poly = NewSimplePolygon([]Point2D{
		{10.884159114175638, -6.233883493017709},
		{15.187183935502128, 1.001890978811065},
		{-14.090321815841484, 6.233904878470094},
		{-15.18718918268496, -1.558457309876673},
		{-2.1937057966535596, -3.896181854342805},
		{10.884159114175638, -6.233883493017709},
	})
	res = poly.ContainsPoint(Point2D{0, 0})
	assert.Equal(t, true, res)

	points := []Point2D{
		{-3.6498172830260374, 4.298256659603395},
		{5.519337229604841, -1.1545251555347096},
		{3.6498120358432047, -4.29823527415101},
		{-5.519342476787673, 1.1545465409870952},
		{-3.6498172830260374, 4.298256659603395},
	}
	res = poly.ContainsAllPoints(points)
	assert.Equal(t, true, res)
}

func TestPolyInPoly(t *testing.T) {
	parcel := NewSimplePolygon(
		[]Point2D{
			{2.194131, 13.191360},
			{-2.194131, 12.857402},
			{-0.506338, 10.853650},
			{0.253170, -13.191360},
		})
	footprint := NewSimplePolygon(
		[]Point2D{
			{-3.6576 / 2, 10.668 / 2},
			{3.6576 / 2, 10.668 / 2},
			{3.6576 / 2, -10.668 / 2},
			{-3.6576 / 2, -10.668 / 2},
		}).
		Rotate(4.052654523130833).
		Shift(Point2D{-0.277551, 16.671435})
	assert.False(t, parcel.Contains(footprint))
}
func TestPointHoles(t *testing.T) {
	poly := Polygon2D{Polygon: []Point2D{
		{0, 0},
		{0, 1},
		{1, 1},
		{1, 0},
	}, Holes: [][]Point2D{[]Point2D{
		{0.25, 0.25},
		{0.25, 0.75},
		{0.75, 0.75},
		{0.75, 0.25},
	}}}
	res := poly.ContainsPoint(Point2D{0.5, 0.5})
	assert.Equal(t, false, res)

	res = poly.ContainsPoint(Point2D{2.5, 2.5})
	assert.Equal(t, false, res)
}

func TestPointInPoly(t *testing.T) {
	p := NewSimplePolygon([]Point2D{
		{-5.946199, 8.126325},
		{0.126515, -7.681045},
		{1.222980, -8.126323},
		{5.946203, 3.116948},
	})
	fmt.Println(p.DebugString())
	assert.True(t, p.ContainsPoint(Point2D{-2.862660, 5.199420}))
	assert.True(t, p.ContainsPoint(Point2D{0.963063, -4.758997}))
}

func TestIntersections(t *testing.T) {

	outer := []Point2D{
		{0, 0},
		{0, 1},
		{1, 1},
		{1, 0},
	}
	inner := []Point2D{
		{0.25, 0.25},
		{0.25, 0.75},
		{0.75, 0.75},
		{0.75, 0.25},
	}

	res := isLineStringInLineString(inner, outer)
	assert.Equal(t, true, res)
}

func TestPolyIntersections(t *testing.T) {
	p := NewSimplePolygon([]Point2D{
		{0, 0},
		{0, 1},
		{1, 1},
		{1, 0},
	})
	p2 := NewSimplePolygon([]Point2D{
		{0.75, 0.75},
		{0.75, 1.75},
		{1.75, 1.75},
		{1.75, 0.75},
	})
	p3 := NewSimplePolygon([]Point2D{
		{2.75, 2.75},
		{2.75, 3.75},
		{3.75, 3.75},
		{3.75, 2.75},
	})
	assert.True(t, p.Intersects(p2))
	assert.True(t, p2.Intersects(p))

	assert.False(t, p.Intersects(p3))
	assert.False(t, p3.Intersects(p))
}
