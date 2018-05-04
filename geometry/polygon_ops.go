package geometry

func (a Bounds) Intersects(b Bounds) bool {
	// https://stackoverflow.com/questions/306316/determine-if-two-rectangles-overlap-each-other
	// RectA.X1 < RectB.X2 && RectA.X2 > RectB.X1 && RectA.Y1 > RectB.Y2 && RectA.Y2 < RectB.Y1
	// Y1 <=> Y2 because order is different
	return (a.MinX < b.MaxX && a.MaxX > b.MinX && a.MaxY > b.MinY && a.MinY < b.MaxY)
}

func (a Polygon2D) Intersects(b Polygon2D) bool {

	// Fast bounds check
	ab := a.Bounds()
	bb := b.Bounds()
	if !ab.Intersects(bb) {
		return false
	}

	// Simple contains
	if a.Contains(b) || b.Contains(a) {
		return true
	}

	// If not completely contained then intersection: check point-based contains
	for _, p := range a.Polygon {
		if b.ContainsPoint(p) {
			return true
		}
	}
	for _, p := range b.Polygon {
		if a.ContainsPoint(p) {
			return true
		}
	}

	// TODO: Handle Holes

	return false
}
