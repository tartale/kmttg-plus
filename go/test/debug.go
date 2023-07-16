package test

import (
	"fmt"
	"os"
	"path"
	"runtime"
)

func DebugDir() (string, error) {
	var (
		file string
		ok   bool
	)
	if _, file, _, ok = runtime.Caller(0); !ok {
		return "", fmt.Errorf("error while trying to locate debug directory")
	}
	dir := path.Dir(path.Dir(file))
	debugDir := path.Join(dir, "debug")
	err := os.MkdirAll(debugDir, os.FileMode(0755))
	if err != nil {
		return "", fmt.Errorf("error while trying to create debug directory: %w", err)
	}

	return debugDir, nil
}
