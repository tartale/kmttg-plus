package config

import (
	"errors"
	"fmt"
	"os"
	"path"
	"reflect"
	"runtime"
	"time"

	"github.com/go-playground/validator/v10"
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

type values struct {
	LogLevel           string        `mapstructure:"KMTTG_LOG_LEVEL" default:"INFO"`
	LogMessages        bool          `mapstructure:"KMTTG_LOG_MESSAGES" default:"false"`
	Port               string        `mapstructure:"KMTTG_PORT" default:"8080"`
	MediaAccessKey     string        `mapstructure:"KMTTG_MEDIA_ACCESS_KEY" default:""`
	Timeout            time.Duration `mapstructure:"KMTTG_TIMEOUT" default:"10s"`
	MaxBackgroundTasks int           `mapstructure:"KMTTG_MAX_BACKGROUND_TASKS" default:"8"`
	WebUIDir           string        `mapstructure:"KMTTG_WEBUI_DIR" default:""`
	TempDir            string        `mapstructure:"KMTTG_TEMP_DIR" default:"${PWD}/tmp" validate:"dir"`
	OutputDir          string        `mapstructure:"KMTTG_OUTPUT_DIR" default:"${PWD}/output" validate:"dir"`
	TivoDecodePath     string        `mapstructure:"KMTTG_TIVODECODE_PATH" default:"${PWD}/tools/tivodecode/tivodecode" validate:"file"`
	ComskipPath        string        `mapstructure:"KMTTG_COMSKIP_PATH" default:"${PWD}/tools/comskip/comskip" validate:"file"`
}

func (v *values) SetDefaults() {
	defaults.SetDefaults(v)
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

	validate := validator.New(validator.WithRequiredStructEnabled())
	return validate.Struct(Values)
}
