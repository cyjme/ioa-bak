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
	DB struct {
		Host               string `mapstructure:"host"`
		Port               string `mapstructure:"port"`
		User               string `mapstructure:"user"`
		Password           string `mapstructure:"password"`
		Name               string `mapstructure:"name"`
		MaxIdleConnections int    `mapstructure:"max_idle_connections"`
		MaxOpenConnections int    `mapstructure:"max_idle_connections"`
	} `mapstructure:"db"`
	Ioa struct {
		Host string `mapstructure:"host"`
		Port string `mapstructure:"port"`
	} `mapstructure:"ioa"`
	Plugins []plugin `mapstructure:"plugins"`
}

func InitConfig() {
	//local
	viper.SetConfigName("config")
	viper.AddConfigPath("./config")
	viper.SetConfigType("yml")
	viper.ReadInConfig()
	//viper.ReadConfig(bytes.NewBufferString(remoteConfig))
	if err := viper.Unmarshal(&Config); err != nil {
		panic(err)
	}
}
