package ops

import (
	"encoding/json"
	"testing"
)

func assertMerge(t *testing.T, old string, new string, res string) {
	oldDict := make(map[string]interface{})
	newDict := make(map[string]interface{})
	e := json.Unmarshal([]byte(old), &oldDict)
	if e != nil {
		t.Error(e)
		return
	}
	e = json.Unmarshal([]byte(new), &newDict)
	if e != nil {
		t.Error(e)
		return
	}
	resDict, e := Merge(oldDict, newDict)
	if e != nil {
		t.Error(e)
		return
	}

	// Result
	marshaled, e := json.Marshal(&resDict)
	if e != nil {
		t.Error(e)
		return
	}
	resStr := string(marshaled)

	// Expected
	expDict := make(map[string]interface{})
	e = json.Unmarshal([]byte(res), &expDict)
	if e != nil {
		t.Error(e)
		return
	}
	exp, e := json.Marshal(&expDict)
	if e != nil {
		t.Error(e)
		return
	}
	expStr := string(exp)

	if expStr != resStr {
		t.Error("Expected '" + expStr + "', got: '" + resStr + "'")
	}
}

func TestEmptyMerge(t *testing.T) {
	assertMerge(t, `{}`, `{}`, `{}`)
}

func TestDisplayIdMerge(t *testing.T) {
	assertMerge(t, `{"displayId": ["123"]}`, `{}`, `{"displayId": ["123"]}`)
	assertMerge(t, `{"displayId": ["123"]}`, `{"displayId": ["1234"]}`, `{"displayId": ["1234","123"]}`)
	assertMerge(t, `{"displayId": []}`, `{"displayId": ["1234"]}`, `{"displayId": ["1234"]}`)
	assertMerge(t, `{"displayId": ["1","2"]}`, `{"displayId": ["3"]}`, `{"displayId": ["3","1","2"]}`)
}

func TestBasicExtrasMerge(t *testing.T) {
	assertMerge(t,
		`{"extras": {"ints":[{"key": "key_1", "value": 123 }]}}`,
		`{}`,
		`{"extras": {"ints":[{"key": "key_1", "value": 123 }]}}`)
	assertMerge(t,
		`{}`,
		`{"extras": {"ints":[{"key": "key_1", "value": 123 }]}}`,
		`{"extras": {"ints":[{"key": "key_1", "value": 123 }]}}`)
	assertMerge(t,
		`{"extras": {"ints":[{"key": "key_1", "value": 124 }]}}`,
		`{"extras": {"ints":[{"key": "key_1", "value": 123 }]}}`,
		`{"extras": {"ints":[{"key": "key_1", "value": 123 }]}}`)
}

func TestRetiredMerge(t *testing.T) {
	// If current retired status is unknown expect it to be active
	assertMerge(t, `{"retired": true}`, `{}`, `{"retired": false}`)
	assertMerge(t, `{"retired": false}`, `{}`, `{"retired": false}`)

	// If current retired status is true - preserve it
	assertMerge(t, `{"retired": false}`, `{"retired": true}`, `{"retired": true}`)
	assertMerge(t, `{"retired": true}`, `{"retired": true}`, `{"retired": true}`)
	assertMerge(t, `{}`, `{"retired": true}`, `{"retired": true}`)

	// If current retired status is false - preserve it
	assertMerge(t, `{"retired": true}`, `{"retired": false}`, `{"retired": false}`)
	assertMerge(t, `{"retired": false}`, `{"retired": false}`, `{"retired": false}`)
	assertMerge(t, `{}`, `{"retired": false}`, `{"retired": false}`)
}

func TestUnknownFieldsMerge(t *testing.T) {
	assertMerge(t, `{"something": ["123"]}`, `{}`, `{"something": ["123"]}`)
	assertMerge(t, `{}`, `{"something": ["123"]}`, `{"something": ["123"]}`)
	assertMerge(t, `{"something": ["12"]}`, `{"something": ["123"]}`, `{"something": ["123"]}`)
}

func TestGeometryMerge(t *testing.T) {
	// Should use old one if not present in new one
	assertMerge(t, `{"geometry":[[[[1,2]]]]}`, `{}`, `{"geometry":[[[[1,2]]]]}`)
	assertMerge(t, `{"geometry":[[[[1,2]]]],"$geometry_src":[[[[1,1]]]]}`, `{}`, `{"$geometry_src":[[[[1,1]]]],"geometry":[[[[1,2]]]]}`)

	// Should drop $geometry_src if geometry not present
	assertMerge(t, `{}`, `{"$geometry_src":[[[[1,1]]]]}`, `{}`)

	// Should forward $geometry_src if geometry is not changed
	assertMerge(t, `{"geometry":[[[[1,2]]]],"$geometry_src":[[[[1,1]]]]}`, `{"geometry":[[[[1,1]]]]}`, `{"$geometry_src":[[[[1,1]]]],"geometry":[[[[1,2]]]]}`)

	// Should clean $geometry_src if geometry changed
	assertMerge(t, `{"geometry":[[[[1,2]]]],"$geometry_src":[[[[1,1]]]]}`, `{"geometry":[[[[1,3]]]]}`, `{"geometry":[[[[1,3]]]]}`)
	assertMerge(t, `{"geometry":[[[[1,2]]]],"$geometry_src":[[[[1,1]]]]}`, `{"geometry":[[[[1,3]]]],"$geometry_src":[[[[1,4]]]]}`, `{"geometry":[[[[1,3]]]],"$geometry_src":[[[[1,4]]]]}`)

	// Should detect optimized geometry and propagate optimzied version
	assertMerge(t, `{"geometry":[[[[1,2]]]],"$geometry_src":[[[[1,1]]]]}`, `{"geometry":[[[[1,1]]]]}`, `{"geometry":[[[[1,2]]]],"$geometry_src":[[[[1,1]]]]}`)
	assertMerge(t, `{"geometry":[[[[1,2]]]],"$geometry_src":[[[[1,1]]]]}`, `{"geometry":[[[[1,1]]]],"$geometry_src":[[[[1,1]]]]}`, `{"geometry":[[[[1,1]]]],"$geometry_src":[[[[1,1]]]]}`)

	// Should preserve optimized when mergine latest with previous
	assertMerge(t, `{"geometry":[[[[1,2]]]]}`, `{"geometry":[[[[1,2]]]],"$geometry_src":[[[[1,2]]]]}`, `{"geometry":[[[[1,2]]]],"$geometry_src":[[[[1,2]]]]}`)

	// Should proritize latest over previous
	assertMerge(t, `{"geometry":[[[[1,1]]]],"$geometry_src":[[[[1,2]]]]}`, `{"geometry":[[[[1,2]]]],"$geometry_src":[[[[1,2]]]]}`, `{"geometry":[[[[1,2]]]],"$geometry_src":[[[[1,2]]]]}`)
	assertMerge(t, `{"geometry":[[[[1,1]]]],"$geometry_src":[[[[1,2]]]]}`, `{"geometry":[[[[1,2]]]],"$geometry_src":[[[[1,2]]]]}`, `{"geometry":[[[[1,2]]]],"$geometry_src":[[[[1,2]]]]}`)
}
