package config

import (
	"errors"
	"reflect"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/mcuadros/go-defaults"
	"github.com/tartale/go/pkg/errorz"
	"github.com/tartale/go/pkg/stringz"
	"github.com/tartale/go/pkg/structs"
)

var Values values

type values struct {
	ComskipPath        string        `mapstructure:"KMTTG_COMSKIP_PATH" default:"${PWD}/tools/comskip/comskip" validate:"file"`
	LogLevel           string        `mapstructure:"KMTTG_LOG_LEVEL" default:"INFO"`
	LogMessages        bool          `mapstructure:"KMTTG_LOG_MESSAGES" default:"false"`
	MaxBackgroundTasks int           `mapstructure:"KMTTG_MAX_BACKGROUND_TASKS" default:"8"`
	MediaAccessKey     string        `mapstructure:"KMTTG_MEDIA_ACCESS_KEY" default:""`
	OutputDir          string        `mapstructure:"KMTTG_OUTPUT_DIR" default:"${PWD}/output" validate:"dir"`
	Port               string        `mapstructure:"KMTTG_PORT" default:"8080"`
	TempDir            string        `mapstructure:"KMTTG_TEMP_DIR" default:"${PWD}/.tmp" validate:"dir"`
	Timeout            time.Duration `mapstructure:"KMTTG_TIMEOUT" default:"10s"`
	TivoDecodePath     string        `mapstructure:"KMTTG_TIVODECODE_PATH" default:"${PWD}/tools/tivodecode/tivodecode" validate:"file"`
	WebUIDir           string        `mapstructure:"KMTTG_WEBUI_DIR" default:""`
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
