package config

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

		err := viper.BindEnv("KMTTG_MEDIA_ACCESS_KEY")
		if err != nil {
			panic(fmt.Errorf("failed to bind environment variable: %w", err))
		}

		err = viper.BindEnv("KMTTG_LOG_LEVEL")
		if err != nil {
			panic(fmt.Errorf("failed to bind environment variable: %w", err))
		}
		viper.SetDefault("KMTTG_LOG_LEVEL", "INFO")

		err = viper.BindEnv("KMTTG_TIMEOUT")
		if err != nil {
			panic(fmt.Errorf("failed to bind environment variable: %w", err))
		}
		viper.SetDefault("KMTTG_TIMEOUT", 10*time.Second)

		viper.AutomaticEnv() // read in environment variables that match

		// If a config file is found, read it in.
		if err := viper.ReadInConfig(); err == nil {
			fmt.Fprintln(os.Stdout, "using config file:", viper.ConfigFileUsed())
		}

		err = viper.Unmarshal(&Values)
		if err != nil {
			panic(fmt.Errorf("failed to read in config: %w", err))
		}
		fmt.Fprintln(os.Stdout, "config loaded", Values)
	})
}
