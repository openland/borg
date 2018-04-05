package geometry

import (
	"math"
)

type Projection struct {
	center Point3D

	cosLon  float64
	sinLon  float64
	cosNLon float64
	sinNLon float64

	cosLat  float64
	sinLat  float64
	cosNLat float64
	sinNLat float64
}

func NewProjection(center PointGeo) *Projection {
	centerProjected := center.ToCartesian()

	// Avoid multiple calculations of the same values
	cosLon := math.Cos(rad(-center.Longitude))
	sinLon := math.Sin(rad(-center.Longitude))
	cosLat := math.Cos(rad(center.Latitude))
	sinLat := math.Sin(rad(center.Latitude))

	cosNLon := math.Cos(rad(center.Longitude))
	sinNLon := math.Sin(rad(center.Longitude))
	cosNLat := math.Cos(rad(-center.Latitude))
	sinNLat := math.Sin(rad(-center.Latitude))

	return &Projection{
		center: centerProjected,

		cosLon: cosLon,
		sinLon: sinLon,
		cosLat: cosLat,
		sinLat: sinLat,

		cosNLon: cosNLon,
		sinNLon: sinNLon,
		cosNLat: cosNLat,
		sinNLat: sinNLat,
	}
}

//
// Point Projection
//

func (point PointGeo) Project(proj *Projection) Point2D {

	//
	// Projecting of point to a projection plane
	//

	// Plane equation:
	// X * centerX + Y * centerY + Z * centerZ = centerX * centerX + centerY * centerY + centerZ * centerZ
	// Line equations:
	// x = t * pointX
	// y = t * pointY
	// z = t * pointZ
	// Intersection
	// t = (centerX * centerX + centerY * centerY + centerZ * centerZ) / (pointX * centerX + pointY * centerY + pointZ * pointZ)

	res := point.ToCartesian()
	t := (proj.center.X*proj.center.X + proj.center.Y*proj.center.Y + proj.center.Z*proj.center.Z) /
		(res.X*proj.center.X + res.Y*proj.center.Y + res.Z*proj.center.Z)

	//
	// Transform Plane to match XY coordinates
	//
	// 1. Scale by t (to reach projection plane)
	// 1. Move center to center of coordinates
	// 2. Rotate around vertical axis (Z)
	// 3. Rotate around horizontal axis (Y)

	res = res.Miltiply(t)
	res = res.Shift(proj.center.Invert())
	res = res.PrecomputedRotateZ(proj.sinLon, proj.cosLon)
	res = res.PrecomputedRotateY(proj.sinLat, proj.cosLat)

	// Ignoring X coordinate as is is aligned with normal to our plane (and should be zero!)
	return Point2D{res.Y, res.Z}
}

func (point Point2D) Unproject(proj *Projection) PointGeo {

	//
	// Rotate coorinates back
	// 1. Rotate around Y
	// 2. Rotate Z
	// 3. Shift to center
	// 4. Scale to length of WORLD_RADIUS
	// 5. Convert back to geo
	//

	return Point3D{0, point.X, point.Y}.
		PrecomputedRotateY(proj.sinNLat, proj.cosNLat).
		PrecomputedRotateZ(proj.sinNLon, proj.cosNLon).
		Shift(proj.center).
		Identity().
		Miltiply(WORLD_RADIUS).
		ToGeo()
}

//
// Polygon Projections
//

func (poly PolygonGeo) Project(proj *Projection) Polygon2D {
	resPoints := make([]Point2D, 0)
	resHoles := make([]LineString2D, 0)

	// Main line string
	for i := 0; i < len(poly.LineStrings[0]); i++ {
		resPoints = append(resPoints, poly.LineStrings[0][i].Project(proj))
	}

	// Holes
	for l := 1; l < len(poly.LineStrings); l++ {
		ls := make([]Point2D, 0)
		for i := 0; i < len(poly.LineStrings[l]); i++ {
			ls = append(ls, poly.LineStrings[l][i].Project(proj))
		}
		resHoles = append(resHoles, ls)
	}

	return Polygon2D{Polygon: resPoints, Holes: resHoles}
}

func (poly Polygon2D) Unproject(proj *Projection) PolygonGeo {
	res := make([]LineStringGeo, 0)

	// Main String
	main := make(LineStringGeo, 0)
	for i := 0; i < len(poly.Polygon); i++ {
		main = append(main, poly.Polygon[0].Unproject(proj))
	}
	res = append(res, main)

	// Holes
	for i := 0; i < len(poly.Holes); i++ {
		hole := make(LineStringGeo, 0)
		for j := 0; j < len(poly.Holes[i]); j++ {
			hole = append(main, poly.Holes[i][j].Unproject(proj))
		}
		res = append(res, hole)

	}

	return PolygonGeo{LineStrings: res}
}

//
// Multipolygon Projections
//

func (multipoly MultipolygonGeo) Project(proj *Projection) Multipolygon2D {
	res := make([]Polygon2D, 0)
	for _, poly := range multipoly.Polygons {
		res = append(res, poly.Project(proj))
	}
	return Multipolygon2D{Polygons: res}
}

func (multipoly Multipolygon2D) Unproject(proj *Projection) MultipolygonGeo {
	res := make([]PolygonGeo, 0)
	for _, poly := range multipoly.Polygons {
		res = append(res, poly.Unproject(proj))
	}
	return MultipolygonGeo{Polygons: res}
}
