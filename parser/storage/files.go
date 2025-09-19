package storage

import (
	"fmt"
	"os"
)

const (
	ROOT = string(".curgo")
)

func afterRoot(fname string) string {
	return fmt.Sprintf("%s/%s", ROOT, fname)
}

func CreateDir(dir string) {
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func ReadFile(fname string) string {
	CreateDir(ROOT)
	path := afterRoot(fname)
	file, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	return string(file)
}

func fileExists(path string) bool {
	_, err := os.Stat(path)

	return err == nil
}

func CreateFile(fname string) string {
	CreateDir(ROOT)
	path := afterRoot(fname)

	if fileExists(path) {
		return path
	}
	file, err := os.Create(path)

	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	file.Close()
	return path
}
