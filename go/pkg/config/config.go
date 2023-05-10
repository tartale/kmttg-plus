package config

import (
	"errors"
	"fmt"
	"os"
	"path"
	"runtime"
)

const certificateFilename = "cdata.p12"

var (
	Values values
)

func CertificatePath() (string, error) {
	var (
		file string
		ok   bool
	)
	if _, file, _, ok = runtime.Caller(0); !ok {
		return "", fmt.Errorf("error while trying to locate certificate file")
	}
	dir := path.Dir(file)
	certificatePath := path.Join(dir, certificateFilename)

	if _, err := os.Stat(certificatePath); errors.Is(err, os.ErrNotExist) {
		return "", fmt.Errorf("certificate file does not exist in expected path: %s", certificatePath)
	}

	return certificatePath, nil
}

type values struct {
	MediaAccessKey string `mapstructure:"KMTTG_MEDIA_ACCESS_KEY"`
}
