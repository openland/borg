package ops

import (
	"testing"
)

func TestSimplify(t *testing.T) {
	coords := [][]float64{
		{-74.00762243660023, 40.70980127030019},
		{-74.00754574127258, 40.70976248032602},
		{-74.00747424759155, 40.70972631999469},
		{-74.00747157803515, 40.70972925486847},
		{-74.00738856093687, 40.70968726535479},
		{-74.00730976505876, 40.70964741697701},
		{-74.00748959017605, 40.70946256265158},
		{-74.00779024860698, 40.709634951687526},
		{-74.00762590743949, 40.709803026169574},
		{-74.00762243660023, 40.70980127030019},
	}
	optimized := OptimizeLine(coords)

	if len(optimized) != 5 {
		t.Error("Expected 5 points, but got: " + string(len(coords)))
	}

	if optimized[0][0] != optimized[len(optimized)-1][0] || optimized[0][1] != optimized[len(optimized)-1][1] {
		t.Error("First and last points doesn't match!")
	}

	t.Log(coords)
	t.Log(optimized)
}
