package utils

import (
	"fmt"
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
	fmt.Println(DebugMultiPolygon(res))
}
