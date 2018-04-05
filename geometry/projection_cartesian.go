package geometry

import "math"

func (src PointGeo) ToCartesian() Point3D {
	centerX := WORLD_RADIUS * math.Cos(rad(src.Latitude)) * math.Cos(rad(src.Longitude))
	centerY := WORLD_RADIUS * math.Cos(rad(src.Latitude)) * math.Sin(rad(src.Longitude))
	centerZ := WORLD_RADIUS * math.Sin(rad(src.Latitude))
	return Point3D{centerX, centerY, centerZ}
}

func (src Point3D) ToGeo() PointGeo {
	latitude := math.Asin(src.Z/WORLD_RADIUS) * 180 / math.Pi
	var longitude float64
	if src.X > 0 {
		longitude = math.Atan(src.Y/src.X) * 180 / math.Pi
	} else if src.Y > 0 {
		longitude = math.Atan(src.Y/src.X)*180/math.Pi + 180
	} else {
		longitude = math.Atan(src.Y/src.X)*180/math.Pi - 180
	}
	return PointGeo{longitude, latitude}
}
