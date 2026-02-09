package test

import (
	"fmt"
	"io"
	"os"
	"path"
	"runtime"

	"github.com/tartale/kmttg-plus/go/pkg/logz"
	"go.uber.org/zap"
)

func GetDebugDir() (string, error) {
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

func MustGetDebugDir() string {
	debugDir, err := GetDebugDir()
	if err != nil {
		panic(err)
	}

	return debugDir
}

func CreateDebugFile(filename string) (*os.File, error) {
	debugDir, err := GetDebugDir()
	if err != nil {
		return nil, err
	}
	debugFile, err := os.Create(path.Join(debugDir, filename))
	if err != nil {
		return nil, err
	}

	return debugFile, nil
}

func MustCreateDebugFile(filename string) *os.File {
	debugFile, err := CreateDebugFile(filename)
	if err != nil {
		panic(err)
	}

	return debugFile
}

func Debug(input io.WriterTo, filename string) {
	if logz.Logger.Level() == zap.DebugLevel {
		file := MustCreateDebugFile(filename)
		defer file.Close()
		input.WriteTo(file)
	}

}
