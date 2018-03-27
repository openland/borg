package utils

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBounds(t *testing.T) {
	bounds := FindBounds([][][][]float64{[][][]float64{[][]float64{{0, 1}, {1, 1}, {1, 0}, {0, 0}}}})
	assert.Equal(t, 0.0, bounds.MinX)
	assert.Equal(t, 0.0, bounds.MinY)
	assert.Equal(t, 1.0, bounds.MaxX)
	assert.Equal(t, 1.0, bounds.MaxY)
}

func TestProject(t *testing.T) {
	polys := [][][][]float64{[][][]float64{[][]float64{
		{-73.996005, 40.722822},
		{-73.995979, 40.722887},
		{-73.996336, 40.72296},
		{-73.996362, 40.722895},
		{-73.996005, 40.722822},
	}}}
	res := ProjectToPlane(polys)
	assert.InEpsilon(t, 13.962607, res[0][0][0][0], 0.00001)
	assert.InEpsilon(t, -7.681031, res[0][0][0][1], 0.00001)
	assert.InEpsilon(t, 16.156113, res[0][0][1][0], 0.00001)
	assert.InEpsilon(t, -0.445260, res[0][0][1][1], 0.00001)
	assert.InEpsilon(t, -13.962578, res[0][0][2][0], 0.00001)
	assert.InEpsilon(t, 7.681058, res[0][0][2][1], 0.00001)
	assert.InEpsilon(t, -16.156111, res[0][0][3][0], 0.00001)
	assert.InEpsilon(t, 0.445296, res[0][0][3][1], 0.00001)
	assert.Equal(t, res[0][0][0][0], res[0][0][4][0])
	assert.Equal(t, res[0][0][0][1], res[0][0][4][1])
}

func TestDistance(t *testing.T) {
	res := MeasureGreatCircleDistance([]float64{-73.996005, 40.722822}, []float64{-73.995979, 40.722887})
	assert.Equal(t, 7.552482132795148, res)
}

func TestPointProjection(t *testing.T) {
	point1 := ProjectToCartesian([]float64{-73.996005, 40.722822})
	point2 := ProjectToCartesian([]float64{-73.995979, 40.722887})
	res := math.Sqrt((point1[0]-point2[0])*(point1[0]-point2[0]) + (point1[1]-point2[1])*(point1[1]-point2[1]) + (point1[2]-point2[2])*(point1[2]-point2[2]))
	assert.Equal(t, 7.5609426673747056, res) // Almost same as test above (should be a little bit bigger)
}
