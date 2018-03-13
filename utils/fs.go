package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
)

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

func SHA256File(filePath string) (result string, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer file.Close()

	hash := sha256.New()
	_, err = io.Copy(hash, file)
	if err != nil {
		return
	}

	result = hex.EncodeToString(hash.Sum(nil))
	return
}
