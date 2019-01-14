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
		Host      string `mapstructure:"host"`
		Port      string `mapstructure:"port"`
		MaxIdleConns int    `mapstructure:"maxIdleConns"`
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
	viper.ReadInConfig()
	//viper.ReadConfig(bytes.NewBufferString(remoteConfig))
	if err := viper.Unmarshal(&Config); err != nil {
		panic(err)
	}

	return Config
}
