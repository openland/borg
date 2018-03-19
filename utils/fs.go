package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
)

func FileExists(name string) bool {
	if _, e := os.Stat(name); e != nil {
		if os.IsNotExist(e) {
			return false
		}
	}
	return true
}

func removeContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
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
		err := removeContents("tmp")
		if err != nil {
			return err
		}
		err = os.Remove("tmp")
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

func CountLines(r *os.File) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}
	_, err := r.Seek(0, 0)
	if err != nil {
		return 0, err
	}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			_, err := r.Seek(0, 0)
			if err != nil {
				return 0, err
			}
			return count, nil

		case err != nil:
			return count, err
		}
	}
}
