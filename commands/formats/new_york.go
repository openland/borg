package formats

import (
	"errors"
	"fmt"

	"github.com/twpayne/go-geom/encoding/geojson"
)

func NewYorkId(feature *geojson.Feature) (string, error) {
	if feature.Properties["BORO"] == nil {
		return "", errors.New("Empty BORO field")
	}
	borough := feature.Properties["BORO"]
	block := int(feature.Properties["BLOCK"].(float64))
	return fmt.Sprintf("%s%05d", borough, block), nil
}
