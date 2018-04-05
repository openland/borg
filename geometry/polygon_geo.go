package geometry

import "math"

// LineStringGeo is as sequence of Geo Points
type LineStringGeo = []PointGeo

// PolygonGeo is a polygon of Geo Coordinates
type PolygonGeo struct {
	LineStrings []LineStringGeo
}

// MultipolygonGeo is a polygon of Geo Coordinates
type MultipolygonGeo struct {
	Polygons []PolygonGeo
}

type BoundsGeo struct {
	MinLongitude float64
	MinLatitude  float64
	MaxLongitude float64
	MaxLatitude  float64
}

func NewSimpleGeoPolygon(points []PointGeo) PolygonGeo {
	return PolygonGeo{LineStrings: [][]PointGeo{points}}
}

//
// Center
//

// Center is a calculation of a center of Polygon
func (polygon PolygonGeo) CenterAverage() PointGeo {
	count := 0
	latitudeS := 0.0
	longitudeS := 0.0
	for i := 0; i < len(polygon.LineStrings[0]); i++ {
		count++
		latitudeS += polygon.LineStrings[0][i].Latitude
		longitudeS += polygon.LineStrings[0][i].Longitude
	}
	return PointGeo{Latitude: latitudeS / float64(count), Longitude: longitudeS / float64(count)}
}

func (poly PolygonGeo) CenterBounds() PointGeo {
	bounds := poly.Bounds()
	return PointGeo{Latitude: (bounds.MinLatitude + bounds.MaxLatitude) / 2, Longitude: (bounds.MaxLongitude + bounds.MinLongitude) / 2}
}

func (poly PolygonGeo) Center() PointGeo {
	return poly.CenterBounds()
}

func (poly PolygonGeo) Bounds() BoundsGeo {
	maxLon := -math.MaxFloat64
	minLon := math.MaxFloat64
	maxLat := -math.MaxFloat64
	minLat := math.MaxFloat64
	for _, circle := range poly.LineStrings {
		for _, point := range circle {
			if point.Latitude > maxLat {
				maxLat = point.Latitude
			}
			if point.Latitude < minLat {
				minLat = point.Latitude
			}
			if point.Longitude > maxLon {
				maxLon = point.Longitude
			}
			if point.Longitude < minLon {
				minLon = point.Longitude
			}
		}
	}
	return BoundsGeo{MinLatitude: minLat, MinLongitude: minLon, MaxLatitude: maxLat, MaxLongitude: maxLon}

}

// Center is a calculation of a center of MultiPolygon
func (polygon MultipolygonGeo) Center() PointGeo {
	count := 0
	latitudeS := 0.0
	longitudeS := 0.0
	for _, poly := range polygon.Polygons {
		for i := 0; i < len(poly.LineStrings[0]); i++ {
			count++
			latitudeS += poly.LineStrings[0][i].Latitude
			longitudeS += poly.LineStrings[0][i].Longitude
		}
	}
	return PointGeo{Latitude: latitudeS / float64(count), Longitude: longitudeS / float64(count)}
}

//
// Area
//

// Ported from https://github.com/mapbox/geojson-area/blob/master/index.js
func measureRingArea(coords LineStringGeo) float64 {
	var p1, p2, p3 PointGeo
	var lowerIndex, middleIndex, upperIndex, i int
	var area float64
	area = 0
	coordsLength := len(coords)

	if coordsLength > 2 {
		for i = 0; i < coordsLength; i++ {
			if i == coordsLength-2 { // i = N-2
				lowerIndex = coordsLength - 2
				middleIndex = coordsLength - 1
				upperIndex = 0
			} else if i == coordsLength-1 { // i = N-1
				lowerIndex = coordsLength - 1
				middleIndex = 0
				upperIndex = 1
			} else { // i = 0 to N-3
				lowerIndex = i
				middleIndex = i + 1
				upperIndex = i + 2
			}
			p1 = coords[lowerIndex]
			p2 = coords[middleIndex]
			p3 = coords[upperIndex]
			area += (rad(p3.Longitude) - rad(p1.Longitude)) * math.Sin(rad(p2.Latitude))
		}

		// wgs84.RADIUS = 6378137
		area = area * WORLD_RADIUS * WORLD_RADIUS / 2
	}

	return area
}

func (polygon PolygonGeo) Area() float64 {
	var res float64
	for _, poly := range polygon.LineStrings {
		res += measureRingArea(poly)
	}
	return res
}

func (polygon MultipolygonGeo) Area() float64 {
	var res float64
	for _, poly := range polygon.Polygons {
		res += poly.Area()
	}
	return res
}

//
// Serialization
//

func NewGeoMultipolygon(polys [][][][]float64) MultipolygonGeo {
	res := make([]PolygonGeo, 0)
	for _, c := range polys {
		res = append(res, NewGeoPolygon(c))
	}
	return MultipolygonGeo{Polygons: res}
}

func NewGeoPolygon(strings [][][]float64) PolygonGeo {
	lineStrings := make([]LineStringGeo, 0)
	for _, s := range strings {
		points := make([]PointGeo, 0)
		for i := 0; i < len(s); i++ {
			points = append(points, PointGeo{Longitude: s[i][0], Latitude: s[i][1]})
		}
		lineStrings = append(lineStrings, points)
	}
	return PolygonGeo{LineStrings: lineStrings}
}

func (poly PolygonGeo) Serialize() [][][]float64 {
	res := make([][][]float64, 0)
	for _, s := range poly.LineStrings {
		points := make([][]float64, 0)
		for i := 0; i < len(s); i++ {
			points = append(points, []float64{s[i].Longitude, s[i].Latitude})
		}
		// Doubling first point
		points = append(points, []float64{s[0].Longitude, s[0].Latitude})
		res = append(res, points)
	}
	return res
}

func (multipolygon MultipolygonGeo) Serialize() [][][][]float64 {
	res := make([][][][]float64, 0)
	for _, p := range multipolygon.Polygons {
		res = append(res, p.Serialize())
	}
	return res
}
