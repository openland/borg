package formats

import (
	"github.com/twpayne/go-geom/encoding/geojson"
)

// IDFunction is a function for ID resolvers
type IDFunction = func(feature *geojson.Feature) (string, error)

// Format handles all configuration for specific format
type Format struct {
	ID IDFunction
}

// Formats provides all available formats
func Formats() map[string]Format {
	res := make(map[string]Format)
	res["nyc_blocks"] = Format{ID: NewYorkBlockID}
	res["nyc_parcels"] = Format{ID: NewYorkParcelID}
	return res
}
