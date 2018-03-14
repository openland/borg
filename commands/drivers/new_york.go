package drivers

import (
	"errors"
	"fmt"

	"github.com/twpayne/go-geom/encoding/geojson"
)

func newYorkBlockID(feature *geojson.Feature) ([]string, error) {
	// Checking fields
	if feature.Properties["BORO"] == nil {
		return []string{}, errors.New("Empty BORO field")
	}
	if feature.Properties["BLOCK"] == nil {
		return []string{}, errors.New("Empty BLOCK field")
	}

	// Basic variables
	borough := feature.Properties["BORO"].(string)
	block := int(feature.Properties["BLOCK"].(float64))

	// Main Format: 400532
	formats := []string{fmt.Sprintf("%s%05d", borough, block)}
	// Secondary Format: 4-532
	formats = append(formats, fmt.Sprintf("%s-%d", borough, block))

	return formats, nil
}

func newYorkParcelID(feature *geojson.Feature) ([]string, error) {
	// Checking fields
	if feature.Properties["BBL"] == nil {
		return []string{}, errors.New("Empty BBL field")
	}

	// Main Format: 4005320024
	formats := []string{feature.Properties["BBL"].(string)}

	// Secondary Format: 4-532-24
	borough := feature.Properties["BORO"].(string)
	block := int(feature.Properties["BLOCK"].(float64))
	lot := int(feature.Properties["LOT"].(float64))
	formats = append(formats, fmt.Sprintf("%s-%d-%d", borough, block, lot))

	return formats, nil
}

// NewYorkBlocksDriver driver for NYC blocks datasets
func NewYorkBlocksDriver() Driver {
	return Driver{ID: newYorkBlockID, Extras: EmptyExtras}
}

// NewYorkParcelsDriver driver for NYC parcels datasets
func NewYorkParcelsDriver() Driver {
	return Driver{ID: newYorkParcelID, Extras: EmptyExtras}
}
