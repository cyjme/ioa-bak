package ioa

import (
	"github.com/spf13/viper"
)

type plugin struct {
	Name string `mapstructure:"name"`
	Path string `mapstructure:"path"`
}

type Config struct {
	Proxy struct {
		Host                string `mapstructure:"host"`
		Port                string `mapstructure:"port"`
		MaxIdleConns        int    `mapstructure:"maxIdleConns"`
		MaxIdleConnsPerHost int    `mapstructure:"maxIdleConnsPerHost"`
	} `mapstructure:"proxy"`
	Plugins []plugin `mapstructure:"plugins"`
}

func ReadConfig() Config {
	var Config Config
	//local
	viper.SetConfigName("config")
	viper.AddConfigPath("./")
	viper.SetConfigType("yml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Error(ERR_CONFIG_LOAD, err)
		panic(err)
	}
	if err := viper.Unmarshal(&Config); err != nil {
		log.Error(ERR_CONFIG_LOAD, err)
		panic(err)
	}

	return Config
}
