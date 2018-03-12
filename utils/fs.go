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

func PrepareTemp() error {
	err := ClearTemp()
	if err != nil {
		return err
	}
	err = os.Mkdir("tmp", os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func ClearTemp() error {
	if FileExists("tmp") {
		err := os.Remove("tmp")
		if err != nil {
			return err
		}
	}
	return nil
}
