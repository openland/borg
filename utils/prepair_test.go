package utils

import "testing"

func TestSplitting(t *testing.T) {
	res := splitBrackets("(a,b,c,d),(e,f,g)")
	if res[0] != "a,b,c,d" {
		t.Errorf("Expected %s, got %s", "a,b,c,d", res[0])
	}
	if res[1] != "e,f,g" {
		t.Errorf("Expected %s, got %s", "e,f,g", res[1])
	}

	res = splitBrackets("(a,b,c,d),(e,f,g),(h,i,j)")
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
