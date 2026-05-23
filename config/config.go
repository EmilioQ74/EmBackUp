package config

import (
	"github.com/spf13/viper"
)

func Init() error {
	viper.SetConfigName(".emBackup")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return viper.SafeWriteConfig()
	}

	return nil
}
