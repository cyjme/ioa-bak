package main

import (
	"encoding/json"
	"errors"
	"ioa"
	"ioa/proto"
	"log"
	"net/http"
	"strings"
)

type ioaPlugin struct {
}

type Data struct {
}

type Config struct {
	Ips string `json:"ips"`
}

type RawConfig struct {
	Ips string `json:"ips"`
}

func (c *Config) UnmarshalJSON(b []byte) error {
	rawConfig := RawConfig{}
	err := json.Unmarshal(b, &rawConfig)
	if err != nil {
		panic(err)
	}
	c.Ips = string(rawConfig.Ips)
	return nil
}

var name = "ip_black"

func (s ioaPlugin) GetName() string {
	return "ip_black"
}

func (s ioaPlugin) GetDescribe() string {
	return "ip_black forbid ip request"
}

func (s ioaPlugin) GetConfigTemplate() proto.ConfigTpl {
	configTpl := proto.ConfigTpl{
		{Name: "ips", Desc: "blackIpList separated by ,", Required: true, FieldType: "string"},
	}

	return configTpl
}

func (i ioaPlugin) InitApi(api *ioa.Api) error {
	err := i.InitApiConfig(api)
	if err != nil {
		return i.throwErr(err)
	}
	err = i.InitApiData(api)
	if err != nil {
		return i.throwErr(err)
	}

	return nil
}

func (i ioaPlugin) InitApiData(api *ioa.Api) error {
	return nil
}

func (i ioaPlugin) InitApiConfig(api *ioa.Api) error {
	var config Config
	json.Unmarshal(api.PluginRawConfig[name], &config)
	log.Println("this is config***********", config)
	api.PluginConfig[name] = config
	return nil
}

func (s ioaPlugin) Run(w http.ResponseWriter, r *http.Request, api *ioa.Api) error {
	addr := r.RemoteAddr
	ip := addr[0:strings.LastIndex(addr, ":")]
	config := api.PluginConfig[name].(Config)
	blackIps := strings.Split(config.Ips, ",")
	log.Println("request ip:", ip)

	for _, i := range blackIps {
		if i == ip {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("request ip is in the ipBlackList"))
			return errors.New("ip forbidden")
		}
	}

	return nil
}

func (i ioaPlugin) throwErr(err error) error {
	return errors.New("plugin" + name + err.Error())
}

var IoaPlugin ioaPlugin