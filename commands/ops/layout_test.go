package ops

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/statecrafthq/borg/utils"
)

func TestLayout(t *testing.T) {
	polys := [][][][]float64{[][][]float64{[][]float64{
		{-73.998824, 40.716576},
		{-73.998862, 40.716515},
		{-73.998523, 40.716396},
		{-73.998485, 40.716457},
		{-73.998824, 40.716576},
	}}}
	proj := utils.NewProjection(polys)
	projected := proj.ProjectMultiPolygon(polys)
	layout := LayoutRectangle(projected, 3.6576, 10.668)
	assert.True(t, layout.Analyzed)
	assert.True(t, layout.Fits)
	assert.True(t, layout.HasLocation)
	assert.InEpsilon(t, 2.0045, layout.Angle, 0.0001)
}
