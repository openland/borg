package ops

import (
	"testing"

	"github.com/statecrafthq/borg/geometry"

	"github.com/stretchr/testify/assert"
)

func TestLayout(t *testing.T) {
	poly := geometry.NewProjectedSimplePolygon([]geometry.PointGeo{
		{-73.998824, 40.716576},
		{-73.998862, 40.716515},
		{-73.998523, 40.716396},
		{-73.998485, 40.716457},
	})
	layout := LayoutRectangle(poly, 3.6576, 10.668)
	assert.True(t, layout.Analyzed)
	assert.True(t, layout.Fits)
	assert.True(t, layout.HasLocation)
	assert.InEpsilon(t, 2.0045, layout.Angle, 0.0001)
}

func TestConvex(t *testing.T) {
	poly := geometry.NewProjectedSimplePolygon([]geometry.PointGeo{
		{-73.999994, 40.72815},
		{-73.999787, 40.728396},
		{-74.000423, 40.728709},
		{-74.000626, 40.728467},
		{-74.000314, 40.728308},
	})
	layout := LayoutRectangle(poly, 3.6576, 10.668)
	assert.True(t, layout.Analyzed)
	assert.True(t, layout.Fits)
	assert.True(t, layout.HasLocation)
}

func TestNarrow(t *testing.T) {
	poly := geometry.NewProjectedSimplePolygon([]geometry.PointGeo{
		{-74.008608, 40.717205},
		{-74.008557, 40.71727},
		{-74.008904, 40.717317},
		{-74.008917, 40.717247},
		{-74.008763, 40.717226},
	})
	layout := LayoutRectangle(poly, 3.6576, 10.668)
	assert.True(t, layout.Analyzed)
	assert.True(t, layout.Fits)
	assert.True(t, layout.HasLocation)
}

func TestComplex(t *testing.T) {
	poly := geometry.NewProjectedSimplePolygon([]geometry.PointGeo{
		{-74.002324, 40.719039},
		{-74.002587, 40.719159},
		{-74.00265, 40.719208},
		{-74.002773, 40.719072},
		{-74.00243, 40.718915},
	})
	layout := LayoutRectangle(poly, 3.6576, 10.668)
	assert.True(t, layout.Analyzed)
	assert.True(t, layout.Fits)
	assert.True(t, layout.HasLocation)
}
