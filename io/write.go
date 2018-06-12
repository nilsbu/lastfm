package io

import (
	"os"
	"path/filepath"
)

func write(data []byte, path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, 0040755); err != nil {
			return err
		}
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}

	_, err = f.Write(data)
	return err
}
