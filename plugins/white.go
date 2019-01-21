package main

import (
	"encoding/json"
	"errors"
	"ioa"
	"ioa/proto"
	"net/http"
	"strings"
)

type Plugin struct {
	ioa.BasePlugin
}

type Data struct {
}

type Config struct {
	Ips []string `json:"ips"`
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

	c.Ips = strings.Split(rawConfig.Ips, ",")

	return nil
}

var name = "ip_white"

func (i Plugin) GetName() string {
	return name
}

func (i Plugin) GetDescribe() string {
	return "ip_white allow ip request"
}

func (i Plugin) GetConfigTemplate() proto.ConfigTpl {
	configTpl := proto.ConfigTpl{
		{Name: "ips", Desc: "whiteIpList separated by , (e.g.: 127.0.0.1,0.0.0.0)", Required: true, FieldType: "string"},
	}

	return configTpl
}

func (i Plugin) InitApi(api *ioa.Api) error {
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

func (i Plugin) InitApiData(api *ioa.Api) error {
	return nil
}

func (i Plugin) InitApiConfig(api *ioa.Api) error {
	var config Config
	json.Unmarshal(api.PluginRawConfig[name], &config)
	api.PluginConfig[name] = config
	return nil
}

func (i Plugin) Run(w http.ResponseWriter, r *http.Request, api *ioa.Api) error {
	addr := r.RemoteAddr
	ip := addr[0:strings.LastIndex(addr, ":")]
	config := api.PluginConfig[name].(Config)
	i.Logger().Debug("request ip:", ip)

	for _, i := range config.Ips {
		if i != ip {
			continue
		}
		return nil
	}
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("request ip is not in the ipWhiteList"))
	return errors.New("ip forbidden")
}

func (i Plugin) throwErr(err error) error {
	return errors.New("plugin" + name + err.Error())
}

var ExportPlugin Plugin
