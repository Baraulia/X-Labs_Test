package main

//nolint:depguard
import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Logger LoggerConf
	GRPC   GRPCConf
}

type LoggerConf struct {
	Level string `mapstructure:"level" default:"INFO"`
}

type GRPCConf struct {
	Port string `mapstructure:"port" default:"50051"`
}

func NewConfig(path string) (Config, error) {
	var conf Config
	viper.SetConfigFile(path)

	if err := viper.ReadInConfig(); err != nil {
		return conf, fmt.Errorf("error while reading config file: %w", err)
	}

	if err := viper.Unmarshal(&conf); err != nil {
		return conf, fmt.Errorf("error while unmarshaling config: %w", err)
	}

	return conf, nil
}
