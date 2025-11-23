package config

import (
	"github.com/gofiber/fiber/v3/log"
	"github.com/spf13/viper"
)

type config struct {
	DSN  string
	Port string
}

var C config

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath("./")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	err = viper.Unmarshal(&C)
	if err != nil {
		log.Fatal(err)
	}
}
