package geometry

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSegment(t *testing.T) {
	segment := Vector2D{Origin: Point2D{X: 1, Y: 1}, DX: 3, DY: 2}
	n := segment.Normal()
	assert.Equal(t, 2.0, n.DX)
	assert.Equal(t, -3.0, n.DY)

	// shift := Vector2D{Origin: n.Identity().Multiply(10).Origin, DX: 3, DY: 2}
}
