package config

import (
	"fmt"
	"github.com/spf13/viper"
	"k8s.io/client-go/util/homedir"
)

type server struct {
	Name string `mapstructure:"name"`
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

var AppConfig = server{}

func init() {
	viper.SetConfigName("application")
	viper.SetConfigType("yml")
	viper.AddConfigPath(homedir.HomeDir())
	viper.AddConfigPath(".")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	viper.UnmarshalKey("server", &AppConfig)
}
