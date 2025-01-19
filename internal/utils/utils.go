package utils

import (
	"errors"
	"io/fs"
	"os"
)

const (
	TIME_FORMAT = "2006-01-01 01:01:01"
)

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}
	return false, err
}

func DeleteFileIfExists(path string) error {
	exists, err := PathExists(path)
	if err != nil {
		return err
	}

	if exists {
		err := os.Remove(path)
		if err != nil {
			return err
		}
	}

	return nil
}

func CreateDirIfNotExists(dirPath string, permissions os.FileMode) error {
	exists, _ := PathExists(dirPath)
	if !exists {
		err := os.MkdirAll(dirPath, permissions)
		if err != nil {
			return err
		}
	}

	return nil
}
