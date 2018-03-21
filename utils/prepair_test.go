package utils

import "testing"

func TestSplitting(t *testing.T) {
	res, e := splitBrackets("(a,b,c,d),(e,f,g)")
	if e != nil {
		t.Error(e)
	}
	if res[0] != "a,b,c,d" {
		t.Errorf("Expected %s, got %s", "a,b,c,d", res[0])
	}
	if res[1] != "e,f,g" {
		t.Errorf("Expected %s, got %s", "e,f,g", res[1])
	}

	res, e = splitBrackets("(a,b,c,d),(e,f,g),(h,i,j)")
	if e != nil {
		t.Error(e)
	}
	if res[0] != "a,b,c,d" {
		t.Errorf("Expected %s, got %s", "a,b,c,d", res[0])
	}
	if res[1] != "e,f,g" {
		t.Errorf("Expected %s, got %s", "e,f,g", res[1])
	}
	if res[2] != "h,i,j" {
		t.Errorf("Expected %s, got %s", "h,i,j", res[2])
	}
}

func TestThinLine(t *testing.T) {
	// [[[[-73.92322762969822 40.828852037931064] [-73.9235617653475 40.828967083658526] [-73.92356166246613 40.82896725925209] [-73.9232275268162 40.828852213524335] [-73.92322762969822 40.828852037931064]]]]
	// TODO: Implement test
}
