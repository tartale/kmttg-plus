package logz

import (
	"encoding"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"runtime"

	"go.uber.org/zap"
)

type FullMarshaler interface {
	encoding.TextMarshaler
	json.Marshaler
}

func InitDebugDir() error {

	if Logger.Level() == zap.DebugLevel {
		debugDir := MustGetDebugDir()
		os.RemoveAll(debugDir)
		err := os.MkdirAll(debugDir, os.FileMode(0755))
		if err != nil {
			return fmt.Errorf("error while trying to create debug directory: %w", err)
		}
	}

	return nil
}

func GetDebugDir() (string, error) {
	var (
		file string
		ok   bool
	)
	if _, file, _, ok = runtime.Caller(0); !ok {
		return "", fmt.Errorf("error while trying to locate debug directory")
	}
	rootDir := path.Dir(file)
	for rootDir != "/" {
		gitFolder := path.Join(rootDir, ".git")
		if _, err := os.Stat(gitFolder); errors.Is(err, os.ErrNotExist) {
			rootDir = path.Dir(rootDir)
			continue
		}
		break
	}
	debugDir := path.Join(rootDir, "debug")

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

func Debug(m FullMarshaler, filename string) {
	if Logger.Level() >= zap.DebugLevel {
		textFile := MustCreateDebugFile(filename + ".txt")
		defer textFile.Close()

		textBytes, err := m.MarshalText()
		if err == nil {
			textFile.Write(textBytes)
		}

		jsonFile := MustCreateDebugFile(filename + ".json")
		defer jsonFile.Close()
		jsonBytes, err := m.MarshalJSON()
		if err == nil {
			jsonFile.Write(jsonBytes)
		}
	}
}
