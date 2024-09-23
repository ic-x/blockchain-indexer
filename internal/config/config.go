package config

import (
	"log"

	"github.com/spf13/viper"
)

func setDefaults() {
	viper.SetDefault("retry_interval", 10)
	viper.SetDefault("block_buffer_size", 0)
	viper.SetDefault("headers_buffer_size", 0)
	viper.SetDefault("out", "blocks.log")
}

func LoadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		log.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		log.Println("No config file found, using default values")
	}

	setDefaults()
}
