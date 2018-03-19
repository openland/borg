package drivers

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/statecrafthq/borg/utils"
)

func newYorkBlockID(feature *utils.Feature) ([]string, error) {
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

func newYorkParcelID(feature *utils.Feature) ([]string, error) {

	// Checking fields
	if feature.Properties["BBL"] == nil {
		return []string{}, errors.New("Empty BBL field")
	}

	// Main ID in format: 4005320024
	var bbl string
	switch v := feature.Properties["BBL"].(type) {
	case string:
		bbl = v
	case float64:
		bbl = strconv.FormatInt(int64(v), 10)
	default:
		return []string{}, errors.New("Unsupported BBL value")
	}
	formats := []string{bbl}

	// Secondary Format: 4-532-24
	var borough string
	if feature.Properties["BORO"] != nil {
		borough = feature.Properties["BORO"].(string)
	} else if feature.Properties["Borough"] != nil {
		boroughKey := strings.ToLower(feature.Properties["Borough"].(string))
		switch boroughKey {
		case "qn":
			borough = "4"
		case "mn":
			borough = "1"
		case "bx":
			borough = "2"
		case "bk":
			borough = "3"
		case "si":
			borough = "5"
		default:
			return []string{}, errors.New("Unknown borough value " + boroughKey)
		}
	} else {
		return formats, nil
	}

	var block int
	if feature.Properties["BLOCK"] != nil {
		block = int(feature.Properties["BLOCK"].(float64))
	} else if feature.Properties["Block"] != nil {
		block = int(feature.Properties["Block"].(float64))
	} else {
		return formats, nil
	}

	var lot int
	if feature.Properties["LOT"] != nil {
		lot = int(feature.Properties["LOT"].(float64))
	} else if feature.Properties["Lot"] != nil {
		lot = int(feature.Properties["Lot"].(float64))
	} else {
		return formats, nil
	}

	formats = append(formats, fmt.Sprintf("%s-%d-%d", borough, block, lot))

	return formats, nil
}

func newYorkParcelExtras(feature *utils.Feature, extras *Extras) error {
	if feature.Properties["ZoneDist1"] != nil ||
		feature.Properties["ZoneDist2"] != nil ||
		feature.Properties["ZoneDist3"] != nil ||
		feature.Properties["ZoneDist4"] != nil {
		zoning := []string{}
		if feature.Properties["ZoneDist1"] != nil {
			zoning = append(zoning, feature.Properties["ZoneDist1"].(string))
		}
		if feature.Properties["ZoneDist2"] != nil {
			zoning = append(zoning, feature.Properties["ZoneDist2"].(string))
		}
		if feature.Properties["ZoneDist3"] != nil {
			zoning = append(zoning, feature.Properties["ZoneDist3"].(string))
		}
		if feature.Properties["ZoneDist4"] != nil {
			zoning = append(zoning, feature.Properties["ZoneDist4"].(string))
		}
		extras.AppendEnum("zoning", zoning)
	}
	if feature.Properties["UnitsTotal"] != nil {
		extras.AppendInt("count_rooms", int32(feature.Properties["UnitsTotal"].(float64)))
	}
	if feature.Properties["NumBldgs"] != nil {
		extras.AppendInt("count_units", int32(feature.Properties["NumBldgs"].(float64)))
	}
	if feature.Properties["NumFloors"] != nil {
		extras.AppendInt("count_stories", int32(feature.Properties["NumFloors"].(float64)))
	}
	if feature.Properties["YearBuilt"] != nil {
		extras.AppendInt("year_built", int32(feature.Properties["YearBuilt"].(float64)))
	}
	if feature.Properties["Address"] != nil {
		extras.AppendString("address", feature.Properties["Address"].(string))
	}
	if feature.Properties["OwnerName"] != nil {
		extras.AppendString("owner_name", feature.Properties["OwnerName"].(string))
	}
	if feature.Properties["AssessLand"] != nil {
		extras.AppendInt("land_value", int32(feature.Properties["AssessLand"].(float64)))
	}
	return nil
}

// NewYorkBlocksDriver driver for NYC blocks datasets
func NewYorkBlocksDriver() Driver {
	return Driver{ID: newYorkBlockID, Extras: EmptyExtras, Record: IgnoreWithoutGeometry}
}

// NewYorkParcelsDriver driver for NYC parcels datasets
func NewYorkParcelsDriver() Driver {
	return Driver{ID: newYorkParcelID, Extras: newYorkParcelExtras, Record: IgnoreWithoutGeometry}
}
