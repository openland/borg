package drivers

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/statecrafthq/borg/commands/ops"
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

func newYorkRecordType(feature *utils.Feature) (RecordType, error) {
	// Check Geometry
	if feature.Geometry == nil {
		return Ignored, nil
	}

	// Check Id
	var bbl string
	switch v := feature.Properties["BBL"].(type) {
	case string:
		bbl = v
	case float64:
		bbl = strconv.FormatInt(int64(v), 10)
	default:
		return Ignored, errors.New("Unsupported BBL value")
	}
	if bbl == "2024760045" || bbl == "2026230616" || bbl == "2028440048" || bbl == "2050100175" || bbl == "2054080138" || bbl == "2057560239" {
		return Ignored, nil
	}
	return Primary, nil
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

	formats = append(formats, fmt.Sprintf("%s-%05d-%04d", borough, block, lot))
	formats = append(formats, fmt.Sprintf("%s-%d-%d", borough, block, lot))

	return formats, nil
}

func newYorkParcelExtras(feature *utils.Feature, extras *ops.Extras) error {
	zoning := []string{}
	if feature.Properties["ZoneDist1"] != nil ||
		feature.Properties["ZoneDist2"] != nil ||
		feature.Properties["ZoneDist3"] != nil ||
		feature.Properties["ZoneDist4"] != nil {
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
		buildings := int32(feature.Properties["NumBldgs"].(float64))
		extras.AppendInt("count_units", buildings)
		if buildings == 0 {
			extras.AppendString("is_vacant", "true")
		} else {
			extras.AppendString("is_vacant", "false")
		}
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
	isPublic := false
	isPublicHandled := false
	if feature.Properties["OwnerName"] != nil {
		isPublicHandled = true
		name := feature.Properties["OwnerName"].(string)
		extras.AppendString("owner_name", name)

		// Simple tokenizer
		nameTokenized := strings.ToLower(name)
		nameTokenized = strings.Replace(nameTokenized, "/", " ", -1)
		nameTokenized = strings.Replace(nameTokenized, "\\", " ", -1)
		nameTokenized = strings.Replace(nameTokenized, "(", " ", -1)
		nameTokenized = strings.Replace(nameTokenized, ")", " ", -1)
		nameTokenized = strings.Replace(nameTokenized, ".", " ", -1)
		for strings.Contains(nameTokenized, "  ") {
			nameTokenized = strings.Replace(nameTokenized, "  ", " ", -1)
		}
		tokens := strings.Split(nameTokenized, " ")

		// Query 1
		if (strings.Contains(nameTokenized, "department of") || strings.Contains(nameTokenized, "dept of")) &&
			!strings.Contains(nameTokenized, "urban") &&
			!strings.Contains(nameTokenized, "hud") &&
			!strings.Contains(nameTokenized, "inc") &&
			!strings.Contains(nameTokenized, "corp") &&
			!strings.Contains(nameTokenized, "llc") {
			extras.AppendString("urbyn_query_1", "true")
			isPublic = true
		} else {
			extras.AppendString("urbyn_query_1", "false")
		}

		// Query 2
		if (strings.Contains(nameTokenized, "city of new york") || strings.Contains(nameTokenized, "city of ny") || strings.Contains(nameTokenized, "city of n y")) &&
			!strings.Contains(nameTokenized, "inc") &&
			!strings.Contains(nameTokenized, "corp") {
			extras.AppendString("urbyn_query_2", "true")
			isPublic = true
		} else {
			extras.AppendString("urbyn_query_2", "false")
		}

		// Query 3
		q3Names := []string{"dsbs", "doe", "hpd", "dsny", "nypd", "dhs", "doh", "dcas", "dep", "ddc", "nypl", "dof", "fdny", "nycedc", "hra"}
		has := false
		for _, t := range tokens {
			for _, q3 := range q3Names {
				if q3 == t {
					has = true
					break
				}
			}
			if has {
				break
			}
		}
		if has {
			extras.AppendString("urbyn_query_3", "true")
			isPublic = true
		} else {
			extras.AppendString("urbyn_query_3", "false")
		}

		// Extra condisions
		q4Names := []string{"dot", "lirr", "nyc", "nys"}
		has = false
		for _, t := range tokens {
			for _, q3 := range q4Names {
				if q3 == t {
					has = true
					break
				}
			}
			if has {
				break
			}
		}
		if has {
			isPublic = true
		}
	}
	if feature.Properties["OwnerType"] != nil {
		isPublicHandled = true
		if feature.Properties["OwnerType"] == "C" {
			extras.AppendString("owner_type", "CITY")
			isPublic = true
		}
		if feature.Properties["OwnerType"] == "M" {
			extras.AppendString("owner_type", "MIXED")
		}
		if feature.Properties["OwnerType"] == "P" {
			extras.AppendString("owner_type", "PRIVATE")
		}

		if feature.Properties["OwnerType"] == "O" {
			extras.AppendString("owner_type", "OTHER")
		}
		if feature.Properties["OwnerType"] == "X" {
			extras.AppendString("owner_type", "EXCLUDED")
		}
	}
	if isPublicHandled {
		if isPublic {
			extras.AppendString("owner_public", "true")
		} else {
			extras.AppendString("owner_public", "false")
		}
	}
	area := 0.0
	if feature.Properties["LotArea"] != nil {
		area = feature.Properties["LotArea"].(float64)
		extras.AppendFloat("assessor_area", utils.SqFeetToMeters((feature.Properties["LotArea"].(float64))))
	}
	if feature.Properties["LotFront"] != nil {
		extras.AppendFloat("assessor_front", utils.FeetToMeters(feature.Properties["LotFront"].(float64)))
	}
	if feature.Properties["LotDepth"] != nil {
		extras.AppendFloat("assessor_depth", utils.FeetToMeters(feature.Properties["LotDepth"].(float64)))
	}
	if feature.Properties["AssessLand"] != nil {
		extras.AppendInt("land_value", int32(feature.Properties["AssessLand"].(float64)))
	}

	// Capacity
	if area > 0 && len(zoning) > 0 {
		maxCapacity := 0
		zoningData := NYCZoning()
		for _, z := range zoning {
			for _, zs := range strings.Split(z, "/") {

				// Normalizing zoning codes
				normalized := zs
				if normalized == "R1" {
					normalized = "R1-1"
				}
				if normalized == "R6" {
					normalized = "R6QH"
				}

				// Calculation of capacity
				if zd, ok := zoningData.Zoning[normalized]; ok {
					if zd.DensityFactor > 0 && zd.MaximumFARNarrow > 0 {

						// Base formula
						capacity := area * zd.MaximumFARNarrow / zd.DensityFactor

						// Round down for .74
						// Round up for .75+
						if (capacity - float64(int(capacity))) <= 0.74 {
							capacity = math.Floor(capacity)
						} else {
							capacity = math.Ceil(capacity)
						}

						// Additional rule: If area < 1700 sf > max units <= 2
						if area < 1700 {
							capacity = math.Min(2, capacity)
						}

						// Save result
						if int(capacity) > maxCapacity {
							maxCapacity = int(capacity)
						}
					}
				}
			}
		}
		extras.AppendInt("unit_capacity", int32(maxCapacity))
	}

	// Borough
	var borough string
	hasBorough := false
	if feature.Properties["BORO"] != nil {
		borough = feature.Properties["BORO"].(string)
		hasBorough = true
	} else if feature.Properties["Borough"] != nil {
		boroughKey := strings.ToLower(feature.Properties["Borough"].(string))
		switch boroughKey {
		case "qn":
			borough = "4"
			hasBorough = true
		case "mn":
			borough = "1"
			hasBorough = true
		case "bx":
			borough = "2"
			hasBorough = true
		case "bk":
			borough = "3"
			hasBorough = true
		case "si":
			borough = "5"
			hasBorough = true
		}
	}
	if hasBorough {
		switch borough {
		case "1":
			extras.AppendInt("borough_id", 1)
			extras.AppendString("borough_name", "Manhattan")
		case "2":
			extras.AppendInt("borough_id", 2)
			extras.AppendString("borough_name", "Bronx")
		case "3":
			extras.AppendInt("borough_id", 3)
			extras.AppendString("borough_name", "Brooklyn")
		case "4":
			extras.AppendInt("borough_id", 4)
			extras.AppendString("borough_name", "Queens")
		case "5":
			extras.AppendInt("borough_id", 5)
			extras.AppendString("borough_name", "Staten Island")
		}
	}

	// BBL
	var bbl string
	switch v := feature.Properties["BBL"].(type) {
	case string:
		bbl = v
	case float64:
		bbl = strconv.FormatInt(int64(v), 10)
	default:
		return errors.New("Unsupported BBL value")
	}
	extras.AppendString("nyc_bbl", bbl)

	return nil
}

// NewYorkBlocksDriver driver for NYC blocks datasets
func NewYorkBlocksDriver() Driver {
	return Driver{ID: newYorkBlockID, Extras: EmptyExtras, Record: IgnoreWithoutGeometry, Retired: NoRetired}
}

// NewYorkParcelsDriver driver for NYC parcels datasets
func NewYorkParcelsDriver() Driver {
	return Driver{ID: newYorkParcelID, Extras: newYorkParcelExtras, Record: newYorkRecordType, Retired: NoRetired}
}
