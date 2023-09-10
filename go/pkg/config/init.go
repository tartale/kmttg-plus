package config

import (
	"fmt"
	"os"
	"reflect"
	"sync"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tartale/go/pkg/structs"
)

// This sync primitive allows configuration to be initialized exactly once
// from either main() or in unit test usage
var InitOnce sync.Once

func init() {
	InitConfig("")
}

// InitConfig reads in config file and ENV variables if set.
func InitConfig(cfgFile string) {
	InitOnce.Do(func() {
		if cfgFile != "" {
			// Use config file from the flag.
			viper.SetConfigFile(cfgFile)
		} else {
			// Find home directory.
			home, err := os.UserHomeDir()
			cobra.CheckErr(err)

			// Search config in home directory with name ".kmttg" (without extension).
			viper.AddConfigPath(home)
			viper.SetConfigType("yaml")
			viper.SetConfigName(".kmttg")
		}

		err := Values.BindEnv()
		if err != nil {
			panic(err)
		}
		viper.AutomaticEnv() // read in environment variables that match

		// If a config file is found, read it in.
		if err := viper.ReadInConfig(); err == nil {
			fmt.Fprintln(os.Stdout, "using config file:", viper.ConfigFileUsed())
		}

		err = Values.SetDefaults()
		if err != nil {
			panic(err)
		}

		err = viper.Unmarshal(&Values)
		if err != nil {
			panic(fmt.Errorf("failed to read in config: %w", err))
		}

		err = Values.ResolveVariables()
		if err != nil {
			panic(err)
		}

		err = Values.Validate()
		if err != nil {
			panic(err)
		}

		fmt.Fprintln(os.Stdout, "config loaded", Values)
	})
}

func (v *values) BindEnv() error {

	err := structs.Walk(&Values, func(sf reflect.StructField, sv reflect.Value) error {

		field := structs.NewField(sf, sv)
		tag := field.TagRoot("mapstructure")
		if tag != "" {
			err := viper.BindEnv(tag)
			if err != nil {
				return fmt.Errorf("failed to bind environment variable: %w", err)
			}
		}

		return nil
	})

	return err
}
