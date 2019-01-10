package app

import (
	"github.com/spf13/viper"
)

var Config config

type plugin struct {
	Name string `mapstructure:"name"`
	Path string `mapstructure:"path"`
}

type config struct {
	Http struct {
		Host string `mapstructure:"host"`
		Port string `mapstructure:"port"`
	} `mapstructure:"http"`
	Ioa struct {
		Host string `mapstructure:"host"`
		Port string `mapstructure:"port"`
	} `mapstructure:"ioa"`
	Plugins []plugin `mapstructure:"plugins"`
}

func InitConfig() {
	//local
	viper.SetConfigName("config")
	viper.AddConfigPath("./")
	viper.SetConfigType("yml")
	viper.ReadInConfig()
	//viper.ReadConfig(bytes.NewBufferString(remoteConfig))
	if err := viper.Unmarshal(&Config); err != nil {
		panic(err)
	}
}
