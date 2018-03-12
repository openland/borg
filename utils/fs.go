package utils

import "os"

func FileExists(name string) bool {
	if _, e := os.Stat(name); e != nil {
		if os.IsNotExist(e) {
			return false
		}
	}
	return true
}
