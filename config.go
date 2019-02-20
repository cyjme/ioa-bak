package ioa

import (
	"github.com/spf13/viper"
)

type PluginConfig struct {
	Name string `mapstructure:"name"`
	Path string `mapstructure:"path"`
}

type EtcdConfig struct {
	Endpoints   []string `mapstructure:"endpoints"`
	DialTimeout int      `mapstructure:"dialTimeout"`
	Username    string   `mapstructure:"username"`
	Password    string   `mapstructure:"password"`
}

type Config struct {
	Etcd       EtcdConfig `mapstructure:"etcd"`

	HttpServer struct {
		Host string `mapstructure:"host"`
		Port string `mapstructure:"port"`
	} `mapstructure:"httpServer"`

	Proxy struct {
		Host                string `mapstructure:"host"`
		Port                string `mapstructure:"port"`
		MaxIdleConns        int    `mapstructure:"maxIdleConns"`
		MaxIdleConnsPerHost int    `mapstructure:"maxIdleConnsPerHost"`
	} `mapstructure:"proxy"`

	Plugins []PluginConfig `mapstructure:"plugins"`
}

func ReadConfig(path string) Config {
	var Config Config
	//local
	viper.SetConfigName("config")
	viper.AddConfigPath(path)
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
