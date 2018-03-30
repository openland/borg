package utils

import (
	"errors"
)

func loadKeys(src map[string]interface{}) []string {
	skeys := make([]string, len(src))
	i := 0
	for k := range src {
		skeys[i] = k
		i++
	}
	return skeys
}

func isGeometryChanged(coords1 []interface{}, coords2 []interface{}) bool {
	if len(coords1) != len(coords2) {
		return true
	}
	for i := range coords1 {
		poly1 := coords1[i].([]interface{})
		poly2 := coords2[i].([]interface{})
		if len(poly1) != len(poly2) {
			return true
		}
		for j := range poly1 {
			line1 := poly1[j].([]interface{})
			line2 := poly2[j].([]interface{})
			if len(line1) != len(line2) {
				return true
			}
			for k := range line1 {
				point1 := line1[k].([]interface{})
				point2 := line2[k].([]interface{})
				if len(point1) != len(point2) {
					return true
				}
				for m := range point1 {
					if point1[m] != point2[m] {
						return true
					}
				}
			}
		}
	}
	return false
}

func isKeywordArrayChanged(src []interface{}, dst []interface{}) bool {
	if len(src) != len(dst) {
		return true
	}
	for _, s := range src {
		found := false
		for _, d := range dst {
			if d == s {
				found = true
				break
			}
		}
		if !found {
			return true
		}
	}
	return false
}

func isExtrasValueChanged(src []interface{}, dst []interface{}) bool {
	for _, s := range src {
		for _, d := range dst {
			a := s.(map[string]interface{})
			b := d.(map[string]interface{})
			if a["key"] == b["key"] {
				return a["value"] != b["value"]
			}
		}
	}
	return false
}

func isExtrasKeywordChanged(src []interface{}, dst []interface{}) bool {
	for _, s := range src {
		for _, d := range dst {
			a := s.(map[string]interface{})
			b := d.(map[string]interface{})
			if a["key"] == b["key"] {
				return isKeywordArrayChanged(a["value"].([]interface{}), b["value"].([]interface{}))
			}
		}
		return true
	}
	return false
}

func IsChanged(src map[string]interface{}, dst map[string]interface{}) (bool, error) {

	supportedFlags := make(map[string]bool)
	supportedFlags["id"] = true
	supportedFlags["geometry"] = true
	supportedFlags["displayId"] = true
	supportedFlags["extras"] = true

	supportedExtras := make(map[string]bool)
	supportedExtras["ints"] = true
	supportedExtras["enums"] = true
	supportedExtras["strings"] = true
	supportedExtras["floats"] = true

	// Loading keys
	skeys := loadKeys(src)
	for _, k := range skeys {
		_, ok := supportedFlags[k]
		if !ok {
			return true, errors.New("Unsupported key " + k)
		}
	}

	// Obvious cases
	if len(src) != len(dst) {
		return true, nil
	}
	if len(src) == 0 {
		return false, nil
	}

	// Check if keys are same in two dictionaries
	for _, k := range skeys {
		_, ok := dst[k]
		if !ok {
			return true, nil
		}
	}

	// Check ID field
	if src["id"] != dst["id"] {
		return true, nil
	}

	// Check geometry
	geom1, ok1 := src["geometry"]
	geom2, ok2 := dst["geometry"]
	if ok1 && ok2 {
		coords1 := geom1.([]interface{})
		coords2 := geom2.([]interface{})
		if isGeometryChanged(coords1, coords2) {
			return true, nil
		}
	}

	// Check display id
	displayID1, ok1 := src["displayId"]
	displayID2, ok2 := dst["displayId"]
	if ok1 && ok2 {
		d1 := displayID1.([]interface{})
		d2 := displayID2.([]interface{})
		if isKeywordArrayChanged(d1, d2) {
			return true, nil
		}
	}

	// Retired flag
	retired1, ok1 := src["retired"]
	retired2, ok2 := dst["retired"]
	if ok1 || ok2 {
		var retired1v bool
		var retired2v bool
		retired1v = false
		retired2v = false
		if ok1 {
			retired1v = retired1.(bool)
		}
		if ok2 {
			retired2v = retired2.(bool)
		}
		if retired1v != retired2v {
			return true, nil
		}
	}

	// Check Extras
	extras1, ok1 := src["extras"]
	extras2, ok2 := dst["extras"]
	if ok1 && ok2 {
		e1 := extras1.(map[string]interface{})
		e2 := extras2.(map[string]interface{})
		e1k := loadKeys(e1)
		e2k := loadKeys(e2)

		// Check if we have some random extras
		for _, k := range e1k {
			_, ok := supportedExtras[k]
			if !ok {
				return true, errors.New("Unsupported extras type " + k)
			}
		}
		for _, k := range e2k {
			_, ok := supportedExtras[k]
			if !ok {
				return true, errors.New("Unsupported extras type " + k)
			}
		}

		// Checking value records
		for _, r := range []string{"ints", "strings", "floats"} {
			ints1, ok1 := e1[r]
			ints2, ok2 := e2[r]
			if ok1 && ok2 {
				if isExtrasValueChanged(ints1.([]interface{}), ints2.([]interface{})) {
					return true, nil
				}
			} else if ok1 && !ok2 {
				if len(ints1.([]interface{})) != 0 {
					return true, nil
				}
			} else if !ok1 && ok2 {
				if len(ints2.([]interface{})) != 0 {
					return true, nil
				}
			}
		}

		// Checking enum values
		enums1, ok1 := e1["enums"]
		enums2, ok2 := e2["enums"]
		if ok1 && ok2 {
			if isExtrasKeywordChanged(enums1.([]interface{}), enums2.([]interface{})) {
				return true, nil
			}
		} else if ok1 && !ok2 {
			if len(enums1.([]interface{})) != 0 {
				return true, nil
			}
		} else if !ok1 && ok2 {
			if len(enums2.([]interface{})) != 0 {
				return true, nil
			}
		}
	}
	return false, nil
}
