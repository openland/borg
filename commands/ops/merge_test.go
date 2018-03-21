package ops

import (
	"encoding/json"
	"testing"
)

func assert(t *testing.T, old string, new string, res string) {
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
	assert(t, `{}`, `{}`, `{}`)
}

func TestDisplayIdMerge(t *testing.T) {
	assert(t, `{"displayId": ["123"]}`, `{}`, `{"displayId": ["123"]}`)
	assert(t, `{"displayId": ["123"]}`, `{"displayId": ["1234"]}`, `{"displayId": ["123", "1234"]}`)
	assert(t, `{"displayId": []}`, `{"displayId": ["1234"]}`, `{"displayId": ["1234"]}`)
}

func TestBasicExtrasMerge(t *testing.T) {
	assert(t,
		`{"extras": {"ints":[{"key": "key_1", "value": 123 }]}}`,
		`{}`,
		`{"extras": {"ints":[{"key": "key_1", "value": 123 }]}}`)
	assert(t,
		`{}`,
		`{"extras": {"ints":[{"key": "key_1", "value": 123 }]}}`,
		`{"extras": {"ints":[{"key": "key_1", "value": 123 }]}}`)
	assert(t,
		`{"extras": {"ints":[{"key": "key_1", "value": 124 }]}}`,
		`{"extras": {"ints":[{"key": "key_1", "value": 123 }]}}`,
		`{"extras": {"ints":[{"key": "key_1", "value": 123 }]}}`)
}

func TestRetiredMerge(t *testing.T) {
	// If current retired status is unknown expect it to be active
	assert(t, `{"retired": true}`, `{}`, `{"retired": false}`)
	assert(t, `{"retired": false}`, `{}`, `{"retired": false}`)

	// If current retired status is true - preserve it
	assert(t, `{"retired": false}`, `{"retired": true}`, `{"retired": true}`)
	assert(t, `{"retired": true}`, `{"retired": true}`, `{"retired": true}`)
	assert(t, `{}`, `{"retired": true}`, `{"retired": true}`)

	// If current retired status is false - preserve it
	assert(t, `{"retired": true}`, `{"retired": false}`, `{"retired": false}`)
	assert(t, `{"retired": false}`, `{"retired": false}`, `{"retired": false}`)
	assert(t, `{}`, `{"retired": false}`, `{"retired": false}`)
}

func TestUnknownFieldsMerge(t *testing.T) {
	assert(t, `{"something": ["123"]}`, `{}`, `{"something": ["123"]}`)
	assert(t, `{}`, `{"something": ["123"]}`, `{"something": ["123"]}`)
	assert(t, `{"something": ["12"]}`, `{"something": ["123"]}`, `{"something": ["123"]}`)
}
