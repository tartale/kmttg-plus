package config

import (
	"errors"
	"fmt"
	"os"
	"path"
	"reflect"
	"runtime"
	"time"

	"github.com/mcuadros/go-defaults"
	"github.com/tartale/go/pkg/errorz"
	"github.com/tartale/go/pkg/stringz"
	"github.com/tartale/go/pkg/structs"
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

type RequiredExePath string
type RequiredDir string

type values struct {
	LogLevel       string          `mapstructure:"KMTTG_LOG_LEVEL" default:"INFO"`
	LogMessages    bool            `mapstructure:"KMTTG_LOG_MESSAGES" default:"false"`
	MediaAccessKey string          `mapstructure:"KMTTG_MEDIA_ACCESS_KEY" default:""`
	Timeout        time.Duration   `mapstructure:"KMTTG_TIMEOUT" default:"10s"`
	WebUIDir       string          `mapstructure:"KMTTG_WEBUI_DIR" default:""`
	OutputDir      RequiredDir     `mapstructure:"KMTTG_OUTPUT_DIR" default:"${PWD}/output"`
	TivoDecodePath RequiredExePath `mapstructure:"KMTTG_TIVODECODE_PATH" default:"${PWD}/tools/tivodecode/tivodecode"`
}

func (v *values) SetDefaults() error {

	defaults.SetDefaults(v)

	return nil
}

func (v *values) ResolveVariables() error {

	err := structs.Walk(&Values, func(sf reflect.StructField, sv reflect.Value) error {

		val := sv.Interface()
		err := stringz.Envsubst(&val)
		if err != nil && errors.Is(err, errorz.ErrInvalidType) {
			return nil
		}
		if err != nil {
			return err
		}
		tsv := sv.Type()
		vval := reflect.ValueOf(val)
		sv.Set(vval.Convert(tsv))

		return nil
	})

	return err
}

func (v *values) Validate() error {

	err := structs.Walk(&Values, func(sf reflect.StructField, sv reflect.Value) error {

		switch val := sv.Interface().(type) {
		case RequiredDir:
			dir := string(val)
			err := os.MkdirAll(dir, os.FileMode(0755))
			if err != nil {
				return fmt.Errorf("error while trying to create directory: %w", err)
			}

		case RequiredExePath:
			path := string(val)
			if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
				return fmt.Errorf("error validating file exists: %w", err)
			}
		}

		return nil
	})

	return err
}
