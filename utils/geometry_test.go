package utils

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAngles(t *testing.T) {
	testPoly := [][]float64{
		{-8.690199761239422, 1.391498730615154},
		{5.231000003567562, -7.068785818689552},
		{8.690202895388282, -1.2801690473393077},
		{-5.230990022569944, 7.068789511576092},
		{-8.690199761239422, 1.391498730615154}}
	res := GetAngles(testPoly)

	// Number of Angles should be equals to number of unique points
	assert.Equal(t, 4, len(res))

	// In this case all angles are about 90 degree
	for _, angl := range res {
		assert.InEpsilon(t, 90, angl, 0.1)
	}

	// Random complex polygon from NYC
	nycParcel := [][][][]float64{[][][]float64{[][]float64{
		{-74.00591, 40.719953},
		{-74.005984, 40.719981},
		{-74.005982, 40.719993},
		{-74.006059, 40.720001},
		{-74.006066, 40.719962},
		{-74.005942, 40.719912},
		{-74.00591, 40.719953},
	}}}
	irregular := NewProjection(nycParcel).ProjectMultiPolygon(nycParcel)

	res = GetAngles(irregular[0][0])

	// One angle should be negative as it is a complex polygon
	assert.InEpsilon(t, -70, res[0], 0.1)

	// [-73.951173,40.703062],[-73.951159,40.70307],[-73.95109,40.703115]]]]
	poly2 := [][][][]float64{[][][]float64{[][]float64{
		{-73.95109, 40.703115},
		{-73.951331, 40.70333},
		{-73.951413, 40.703276},
		{-73.951173, 40.703062},
		{-73.951159, 40.70307},
		{-73.95109, 40.703115},
	}}}

	proj := NewProjection(poly2)
	poly3 := proj.ProjectMultiPolygon(poly2)
	res = GetAngles(poly3[0][0])
	fmt.Println(res)
}

func TestSides(t *testing.T) {
	testPoly := [][]float64{
		{-8.690199761239422, 1.391498730615154},
		{5.231000003567562, -7.068785818689552},
		{8.690202895388282, -1.2801690473393077},
		{-5.230990022569944, 7.068789511576092},
		{-8.690199761239422, 1.391498730615154}}
	res := GetSides(testPoly)
	assert.InEpsilon(t, 16.2903, res[0], 0.0001)
	assert.InEpsilon(t, 6.7434, res[1], 0.0001)
	assert.InEpsilon(t, 16.2328, res[2], 0.0001)
	assert.InEpsilon(t, 6.6481, res[3], 0.0001)
}

func TestGlobalAngle(t *testing.T) {

	// 45 degres angle
	angle := GlobalAngle([]float64{0, 0}, []float64{1, 1})
	assert.InEpsilon(t, math.Pi/4, angle, 0.000001)
	angle = GlobalAngle([]float64{1, 1}, []float64{2, 2})
	assert.InEpsilon(t, math.Pi/4, angle, 0.000001)

	// -45 degres angle
	angle = GlobalAngle([]float64{0, 0}, []float64{-1, 1})
	assert.InEpsilon(t, -math.Pi/4, angle, 0.000001)
	angle = GlobalAngle([]float64{1, 1}, []float64{0, 2})
	assert.InEpsilon(t, -math.Pi/4, angle, 0.000001)
}

func TestPointTest(t *testing.T) {
	// WP notes a ray passing through a "side" vertex is an interesting test case.
	// The test case selected for ExampleXY shows the function working properly
	// in this case.
	res := IsPointInside([]float64{1, 1}, [][]float64{{0, 0}, {0, 2}, {2, 1}})
	assert.Equal(t, true, res)

	// Some random polygon
	poly := [][]float64{{17.926384, -31.113776}, {35.388660, -3.729118}, {-18.263668, 31.113820}, {-35.388622, 4.174565}, {-9.068620, -13.525313}}
	res = IsPointInside([]float64{0, 0}, poly)
	assert.Equal(t, true, res)

	//[-15.18718918268496 -1.558457309876673] [-2.1937057966535596 -3.896181854342805] [10.884159114175638 -6.233883493017709]]]]
	poly = [][]float64{
		{10.884159114175638, -6.233883493017709},
		{15.187183935502128, 1.001890978811065},
		{-14.090321815841484, 6.233904878470094},
		{-15.18718918268496, -1.558457309876673},
		{-2.1937057966535596, -3.896181854342805},
		{10.884159114175638, -6.233883493017709},
	}
	fmt.Println(DebugPointLine(poly))
	res = IsPointInside([]float64{0, 0}, poly)
	assert.Equal(t, true, res)

	points := [][]float64{
		{-3.6498172830260374, 4.298256659603395},
		{5.519337229604841, -1.1545251555347096},
		{3.6498120358432047, -4.29823527415101},
		{-5.519342476787673, 1.1545465409870952},
		{-3.6498172830260374, 4.298256659603395},
	}
	res = IsPointsInside(points, poly)
	assert.Equal(t, true, res)
}
