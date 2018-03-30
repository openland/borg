package utils

import (
	"encoding/json"
	"testing"
)

func assertChanged(t *testing.T, src map[string]interface{}, dst map[string]interface{}) {
	r, err := IsChanged(src, dst)
	if err != nil {
		t.Error(err)
		return
	}
	if !r {
		t.Error("Record is not changed while it should")
		return
	}
	r, err = IsChanged(dst, src)
	if err != nil {
		t.Error(err)
		return
	}
	if !r {
		t.Error("Transitive property violated")
	}
}

func assertChangedJson(t *testing.T, src string, dst string) {
	srcDict := make(map[string]interface{})
	dstDict := make(map[string]interface{})
	e := json.Unmarshal([]byte(src), &srcDict)
	if e != nil {
		t.Error(e)
		return
	}
	e = json.Unmarshal([]byte(dst), &dstDict)
	if e != nil {
		t.Error(e)
		return
	}
	assertChanged(t, srcDict, dstDict)
}

func assertNotChanged(t *testing.T, src map[string]interface{}, dst map[string]interface{}) {
	r, err := IsChanged(src, dst)
	if err != nil {
		t.Error(err)
		return
	}
	if r {
		t.Error("Record changed while it should not")
		return
	}
	r, err = IsChanged(dst, src)
	if err != nil {
		t.Error(err)
		return
	}
	if r {
		t.Error("Transitive property violated")
	}
}

func assertNotChangedJson(t *testing.T, src string, dst string) {
	srcDict := make(map[string]interface{})
	dstDict := make(map[string]interface{})
	e := json.Unmarshal([]byte(src), &srcDict)
	if e != nil {
		t.Error(e)
		return
	}
	e = json.Unmarshal([]byte(dst), &dstDict)
	if e != nil {
		t.Error(e)
		return
	}
	assertNotChanged(t, srcDict, dstDict)
}

func TestEmptyDiff(t *testing.T) {
	src := make(map[string]interface{})
	dst := make(map[string]interface{})
	assertNotChanged(t, src, dst)
}

func TestIDField(t *testing.T) {
	src := make(map[string]interface{})
	dst := make(map[string]interface{})

	// Same ID
	src["id"] = "123"
	dst["id"] = "123"
	assertNotChanged(t, src, dst)

	// ID changed
	dst["id"] = "1235"
	assertChanged(t, src, dst)
}

func TestGeometryField(t *testing.T) {
	// Cases
	geomEmpty := `{"geometry":[]}`
	geomSimple1 := `{"geometry":[[[[-74.0176548327901,40.6223216667695],[-74.0176217731316,40.6223017106191],[-74.0178448276858,40.6220835749904],[-74.0179085623536,40.6221220502963],[-74.0176855085386,40.6223401855451],[-74.0176548327901,40.6223216667695]]]]}`
	geomSimple2 := `{"geometry":[[[[-73.9633318910639,40.6229214329515],[-73.9636934930431,40.6228816505593],[-73.9637107745925,40.6229733172174],[-73.9633491727866,40.6230130996637],[-73.9633318910639,40.6229214329515]]]]}`

	// Empty geometry
	assertNotChangedJson(t, geomEmpty, geomEmpty)
	assertNotChangedJson(t, geomSimple1, geomSimple1)
	assertChangedJson(t, geomSimple1, geomSimple2)
}

func TestDisplayIdField(t *testing.T) {
	assertNotChangedJson(t, `{"displayId": ["123", "11"]}`, `{"displayId": ["123", "11"]}`)
	assertChangedJson(t, `{"displayId": ["123"]}`, `{"displayId": ["123", "11"]}`)
	assertChangedJson(t, `{"displayId": ["123"]}`, `{}`)
}

func TestExtras(t *testing.T) {

	// ints
	assertNotChangedJson(t, `{"extras":{"ints":[{"key":"some_key","value":2}]}}`, `{"extras":{"ints":[{"key":"some_key","value":2}]}}`)
	assertChangedJson(t, `{"extras":{"ints":[{"key":"some_key","value":2}]}}`, `{"extras":{"ints":[{"key":"some_key","value":1}]}}`)

	// strings
	assertNotChangedJson(t, `{"extras":{"strings":[{"key":"some_key","value":"957 80 STREET"}]}}`, `{"extras":{"strings":[{"key":"some_key","value":"957 80 STREET"}]}}`)
	assertChangedJson(t, `{"extras":{"strings":[{"key":"some_key","value":"957 80 STREET"}]}}`, `{"extras":{"strings":[{"key":"some_key","value":"957 81 STREET"}]}}`)

	// If not first key is changed (there was a bug)
	assertChangedJson(t, `{"extras":{"strings":[{"key":"some","value":"!"},{"key":"some_key","value":"957 80 STREET"}]}}`, `{"extras":{"strings":[{"key":"some","value":"!"},{"key":"some_key","value":"957 81 STREET"}]}}`)

	// Added key
	assertChangedJson(t, `{"extras":{"strings":[{"key":"some","value":"!"}]}}`, `{"extras":{"strings":[{"key":"some","value":"!"},{"key":"some_key","value":"957 81 STREET"}]}}`)
}

func TestRandomSamples(t *testing.T) {
	assertNotChangedJson(t,
		`{"displayId":["3-5983-66"],"extras":{"enums":[{"key":"zoning","value":["R4-1"]}],"strings":[{"key":"address","value":"957 80 STREET"},{"key":"owner_name","value":"SHERMAN JOHN P"}],"floats":[],"ints":[{"key":"count_rooms","value":1},{"key":"count_units","value":2},{"key":"count_stories","value":2},{"key":"year_built","value":1920},{"key":"land_value","value":13279}]},"geometry":[[[[-74.0176548327901,40.6223216667695],[-74.0176217731316,40.6223017106191],[-74.0178448276858,40.6220835749904],[-74.0179085623536,40.6221220502963],[-74.0176855085386,40.6223401855451],[-74.0176548327901,40.6223216667695]]]],"id":"3059830066"}`,
		`{"displayId":["3-5983-66"],"extras":{"enums":[{"key":"zoning","value":["R4-1"]}],"strings":[{"key":"address","value":"957 80 STREET"},{"key":"owner_name","value":"SHERMAN JOHN P"}],"floats":[],"ints":[{"key":"count_rooms","value":1},{"key":"count_units","value":2},{"key":"count_stories","value":2},{"key":"year_built","value":1920},{"key":"land_value","value":13279}]},"geometry":[[[[-74.0176548327901,40.6223216667695],[-74.0176217731316,40.6223017106191],[-74.0178448276858,40.6220835749904],[-74.0179085623536,40.6221220502963],[-74.0176855085386,40.6223401855451],[-74.0176548327901,40.6223216667695]]]],"id":"3059830066"}`)

	assertChangedJson(t,
		`{"displayId":["1-1-201"],"extras":{"enums":[{"key":"zoning","value":["R3-2"]}],"strings":[{"key":"address","value":"1 ELLIS ISLAND"},{"key":"owner_name","value":"U S GOVT LAND \u0026 BLDGS"},{"key":"owner_type","value":"EXCLUDED"},{"key":"shape_type","value":"convex"},{"key":"analyzed","value":"false"},{"key":"project_kassita1","value":"false"},{"key":"project_kassita2","value":"false"}],"floats":[{"key":"area","value":187490.82414773502}],"ints":[{"key":"count_rooms","value":0},{"key":"count_units","value":8},{"key":"count_stories","value":0},{"key":"year_built","value":1900},{"key":"land_value","value":14972400}]},"geometry":[[[[-74.040028,40.700851],[-74.040925,40.700574],[-74.0451,40.697548],[-74.042371,40.695367],[-74.037543,40.698866],[-74.040028,40.700851]]]],"id":"1000010201","retired":false}`,
		`{"displayId":["1-1-201"],"extras":{"enums":[{"key":"zoning","value":["R3-2"]}],"strings":[{"key":"address","value":"1 ELLIS ISLAND"},{"key":"owner_name","value":"U S GOVT LAND \u0026 BLDGS"},{"key":"owner_type","value":"EXCLUDED"},{"key":"shape_type","value":"convex"},{"key":"analyzed","value":"true"},{"key":"project_kassita1","value":"true"},{"key":"project_kassita2","value":"true"}],"floats":[{"key":"area","value":187490.82414773502},{"key":"project_kassita1_angle","value":-1.9576147956218164},{"key":"project_kassita1_lon","value":-74.04132153737835},{"key":"project_kassita1_lat","value":40.69810900453765},{"key":"project_kassita2_angle","value":-1.9576147956218164},{"key":"project_kassita2_lon","value":-74.04132153737835},{"key":"project_kassita2_lat","value":40.69810900453765}],"ints":[{"key":"count_rooms","value":0},{"key":"count_units","value":8},{"key":"count_stories","value":0},{"key":"year_built","value":1900},{"key":"land_value","value":14972400}]},"geometry":[[[[-74.040028,40.700851],[-74.040925,40.700574],[-74.0451,40.697548],[-74.042371,40.695367],[-74.037543,40.698866],[-74.040028,40.700851]]]],"id":"1000010201","retired":false}`)
}
