package utils

import (
	"io"
	"os"
)

func SaveFile(file io.Reader, path string) (string, error) {
	out, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		return "", err
	}

	return path, nil
}
