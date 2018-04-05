package geometry

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

// func TestBounds(t *testing.T) {
// 	bounds := FindBounds([][][][]float64{[][][]float64{[][]float64{{0, 1}, {1, 1}, {1, 0}, {0, 0}}}})
// 	assert.Equal(t, 0.0, bounds.MinX)
// 	assert.Equal(t, 0.0, bounds.MinY)
// 	assert.Equal(t, 1.0, bounds.MaxX)
// 	assert.Equal(t, 1.0, bounds.MaxY)
// }

func TestMultipolyProject(t *testing.T) {
	polys := NewSimpleGeoPolygon([]PointGeo{
		{-73.996005, 40.722822},
		{-73.995979, 40.722887},
		{-73.996336, 40.72296},
		{-73.996362, 40.722895},
	})
	proj := NewProjection(polys.Center())
	res := polys.Project(proj)
	assert.InEpsilon(t, 13.962607, res.Polygon[0].X, 0.00001)
	assert.InEpsilon(t, -7.681031, res.Polygon[0].Y, 0.00001)
	assert.InEpsilon(t, 16.156113, res.Polygon[1].X, 0.00001)
	assert.InEpsilon(t, -0.445260, res.Polygon[1].Y, 0.00001)
	assert.InEpsilon(t, -13.962578, res.Polygon[2].X, 0.00001)
	assert.InEpsilon(t, 7.681058, res.Polygon[2].Y, 0.00001)
	assert.InEpsilon(t, -16.156111, res.Polygon[3].X, 0.00001)
	assert.InEpsilon(t, 0.445296, res.Polygon[3].Y, 0.00001)

	unproj := res.Unproject(proj)
	for i := 0; i < 4; i++ {
		assert.InEpsilon(t, polys.LineStrings[0][i].Latitude, unproj.LineStrings[0][i].Latitude, 0.00001)
		assert.InEpsilon(t, polys.LineStrings[0][i].Longitude, unproj.LineStrings[0][i].Longitude, 0.00001)
	}
}

func TestPointProject(t *testing.T) {
	polys := NewSimpleGeoPolygon([]PointGeo{
		{-73.996005, 40.722822},
		{-73.995979, 40.722887},
		{-73.996336, 40.72296},
		{-73.996362, 40.722895},
		{-73.996005, 40.722822},
	})
	proj := NewProjection(polys.Center())
	point := PointGeo{-73.996005, 40.722822}.Project(proj)
	assert.InEpsilon(t, 13.962607, point.X, 0.00001)
	assert.InEpsilon(t, -7.681031, point.Y, 0.00001)

	unp := point.Unproject(proj)
	assert.InEpsilon(t, -73.996005, unp.Longitude, 0.000005)
	assert.InEpsilon(t, 40.722822, unp.Latitude, 0.000005)
}

func TestDistance(t *testing.T) {
	res := PointGeo{-73.996005, 40.722822}.DistanceTo(PointGeo{-73.995979, 40.722887})
	assert.Equal(t, 7.552482132795148, res)
}

func TestPointProjection(t *testing.T) {

	// Forward test (should preserve almost same distance as GreatCircleDistance)
	point1 := PointGeo{-73.996005, 40.722822}.ToCartesian()
	point2 := PointGeo{-73.995979, 40.722887}.ToCartesian()
	res := math.Sqrt((point1.X-point2.X)*(point1.X-point2.X) + (point1.Y-point2.Y)*(point1.Y-point2.Y) + (point1.Z-point2.Z)*(point1.Z-point2.Z))
	assert.Equal(t, 7.5609426673747056, res) // Almost same as test above (should be a little bit bigger)

	// Check reverse conversion
	point1u := point1.ToGeo()
	point2u := point2.ToGeo()
	assert.InEpsilon(t, -73.996005, point1u.Longitude, 0.000001)
	assert.InEpsilon(t, 40.722822, point1u.Latitude, 0.000001)
	assert.InEpsilon(t, -73.995979, point2u.Longitude, 0.000001)
	assert.InEpsilon(t, 40.722887, point2u.Latitude, 0.000001)
}
