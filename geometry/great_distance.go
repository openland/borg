package geometry

import (
	"github.com/umahmood/haversine"
)

func (start PointGeo) DistanceTo(to PointGeo) float64 {
	oxford := haversine.Coord{Lat: start.Latitude, Lon: start.Longitude}
	turin := haversine.Coord{Lat: to.Latitude, Lon: to.Longitude}
	_, km := haversine.Distance(oxford, turin)
	return km * 1000
}
