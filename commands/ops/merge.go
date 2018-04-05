package ops

import (
	"encoding/json"

	"github.com/statecrafthq/borg/utils"
)

func mergeDisplayIds(a []interface{}, b []interface{}) []interface{} {
	res := make([]interface{}, 0)

	// Initial filling
	for _, i := range b {
		res = append(res, i)
	}

	// Adding missing in a but present in b
	for _, i := range a {
		found := false
		for _, j := range b {
			if i == j {
				found = true
				break
			}
		}
		if !found {
			res = append(res, i)
		}
	}
	return res
}

func mapHasKey(extras map[string]interface{}, key string) bool {
	_, ok := extras[key]
	if ok {
		return true
	}
	return false
}

func hasKey(extras map[string]interface{}, key string) bool {
	vals, ok := extras["floats"]
	if ok {
		if mapHasKey(vals.(map[string]interface{}), key) {
			return true
		}
	}

	vals, ok = extras["ints"]
	if ok {
		if mapHasKey(vals.(map[string]interface{}), key) {
			return true
		}
	}

	vals, ok = extras["strings"]
	if ok {
		if mapHasKey(vals.(map[string]interface{}), key) {
			return true
		}
	}

	vals, ok = extras["enums"]
	if ok {
		if mapHasKey(vals.(map[string]interface{}), key) {
			return true
		}
	}

	return false
}

func mergeExtrasType(a interface{}, b interface{}) []interface{} {
	res := make([]interface{}, 0)
	added := make(map[string]bool)
	if b != nil {
		bArrr := b.([]interface{})
		for _, v := range bArrr {
			r := v.(map[string]interface{})
			key := r["key"].(string)
			_, present := added[key]
			if !present {
				added[key] = true
				res = append(res, v)
			}
		}
	}
	if a != nil {
		aArrr := a.([]interface{})
		for _, v := range aArrr {
			r := v.(map[string]interface{})
			key := r["key"].(string)
			_, present := added[key]
			if !present {
				added[key] = true
				res = append(res, v)
			}
		}
	}
	return res
}

func MergeExtras(a map[string]interface{}, b map[string]interface{}) (map[string]interface{}, error) {
	res := make(map[string]interface{})
	if a["floats"] != nil || b["floats"] != nil {
		res["floats"] = mergeExtrasType(a["floats"], b["floats"])
	}
	if a["ints"] != nil || b["ints"] != nil {
		res["ints"] = mergeExtrasType(a["ints"], b["ints"])
	}
	if a["strings"] != nil || b["strings"] != nil {
		res["strings"] = mergeExtrasType(a["strings"], b["strings"])
	}
	if a["enums"] != nil || b["enums"] != nil {
		res["enums"] = mergeExtrasType(a["enums"], b["enums"])
	}
	return res, nil
}

func cloneMap(src map[string]interface{}) (map[string]interface{}, error) {
	t, e := json.Marshal(src)
	if e != nil {
		return nil, e
	}
	var res map[string]interface{}
	e = json.Unmarshal(t, &res)
	if e != nil {
		return nil, e
	}
	return res, nil
}

func Merge(previous map[string]interface{}, latest map[string]interface{}) (map[string]interface{}, error) {

	// Cloning
	res, e := cloneMap(latest)
	if e != nil {
		return nil, e
	}
	delete(res, "$geometry_src") // We will forward it manually later

	// Display Id
	ids1, ok1 := previous["displayId"]
	ids2, ok2 := latest["displayId"]
	if ok1 {
		if ok2 {
			res["displayId"] = mergeDisplayIds(ids1.([]interface{}), ids2.([]interface{}))
		} else {
			res["displayId"] = ids1
		}
	}

	// Extras
	ex1, ok1 := previous["extras"]
	ex2, ok2 := latest["extras"]
	if ok1 {
		if ok2 {
			r, e := MergeExtras(ex1.(map[string]interface{}), ex2.(map[string]interface{}))
			if e != nil {
				return nil, e
			}
			res["extras"] = r
		} else {
			res["extras"] = ex1
		}
	}

	// Retired
	_, ok1 = previous["retired"]
	_, ok2 = latest["retired"]
	// If field is missing in latest and present in previous - forward it
	if ok1 && !ok2 {
		res["retired"] = false
	}

	// Geometry
	geometry1, ok1 := previous["geometry"]
	geometry2, ok2 := latest["geometry"]
	goemetry1Src, ok1Src := previous["$geometry_src"]
	goemetry2Src, ok2Src := latest["$geometry_src"]
	if ok1 {
		if ok2 {

			// Detecting real geometry
			realGeometry1 := geometry1
			realGeometry2 := geometry2
			if ok1Src {
				realGeometry1 = goemetry1Src
			}
			if ok2Src {
				realGeometry2 = goemetry2Src
			}

			if utils.IsGeometryChanged(realGeometry1.([]interface{}), realGeometry2.([]interface{})) {
				// Geometry was changed copy from latest
				res["geometry"] = geometry2
				if ok2Src {
					res["$geometry_src"] = goemetry2Src
				}
			} else {
				// Otherwise use fields from latest, if not found copy from previous
				if ok2Src {
					res["geometry"] = geometry2
					res["$geometry_src"] = goemetry2Src
				} else if ok1Src {
					res["geometry"] = geometry1
					res["$geometry_src"] = goemetry1Src
				} else {
					res["geometry"] = geometry2
				}
			}
		} else {
			// Forward geometry
			res["geometry"] = geometry1

			// Forward $geometry_src if present
			goemetry1Src, ok := previous["$geometry_src"]
			if ok {
				res["$geometry_src"] = goemetry1Src
			}
		}
	}

	// Missing keys
	for k := range previous {
		if k == "extras" || k == "displayId" || k == "retired" || k == "geometry" || k == "$geometry_src" {
			continue
		}
		_, p := latest[k]
		if !p {
			res[k] = previous[k]
		}
	}

	return res, nil
}
