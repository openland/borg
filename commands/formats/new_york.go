package formats

import (
	"errors"
	"fmt"

	"github.com/twpayne/go-geom/encoding/geojson"
)

// NewYorkBlockID is a function that builds unique Block ID for NewYork datasets
func NewYorkBlockID(feature *geojson.Feature) (string, error) {
	if feature.Properties["BORO"] == nil {
		return "", errors.New("Empty BORO field")
	}
	borough := feature.Properties["BORO"]
	block := int(feature.Properties["BLOCK"].(float64))
	return fmt.Sprintf("%s%05d", borough, block), nil
}

// NewYorkParcelID is a function that builds unique Parcel ID for NewYork datasets
func NewYorkParcelID(feature *geojson.Feature) (string, error) {
	return feature.Properties["BBL"].(string), nil
}
