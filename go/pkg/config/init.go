package config

import (
	"fmt"
	"os"
	"sync"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tartale/kmttg-plus/go/pkg/logz"
	"go.uber.org/zap"
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
			logz.Logger.Error("failed to bind environment variable", zap.Error(err))
		}

		viper.AutomaticEnv() // read in environment variables that match

		// If a config file is found, read it in.
		if err := viper.ReadInConfig(); err == nil {
			fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		}

		err = viper.Unmarshal(&Values)
		if err != nil {
			logz.Logger.Panic("failed to read in config", zap.Error(err))
		}
		logz.Logger.Info("config loaded", zap.Any("config", Values))
	})
}
