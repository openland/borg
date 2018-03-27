package ops

import (
	"encoding/json"
)

func mergeDisplayIds(a []interface{}, b []interface{}) []interface{} {
	res := make([]interface{}, 0)

	// Initial filling
	for _, i := range a {
		res = append(res, i)
	}

	// Adding missing in a but present in b
	for _, i := range b {
		found := false
		for _, j := range a {
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

	// Missing keys
	for k := range previous {
		if k == "extras" || k == "displayId" || k == "retired" {
			continue
		}
		_, p := latest[k]
		if !p {
			res[k] = previous[k]
		}
	}

	return res, nil
}
