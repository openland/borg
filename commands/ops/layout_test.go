package ops

import (
	"fmt"
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

func TestConvex(t *testing.T) {
	// [[[[-73.999994,40.72815],[-73.999787,40.728396],[-74.000423,40.728709],[-74.000626,40.728467],[-74.000314,40.728308],[-73.999994,40.72815]]]]
	polys := [][][][]float64{[][][]float64{[][]float64{
		{-73.999994, 40.72815},
		{-73.999787, 40.728396},
		{-74.000423, 40.728709},
		{-74.000626, 40.728467},
		{-74.000314, 40.728308},
		{-73.999994, 40.72815},
	}}}
	proj := utils.NewProjection(polys)
	projected := proj.ProjectMultiPolygon(polys)
	layout := LayoutRectangle(projected, 3.6576, 10.668)
	fmt.Println(layout)
}
