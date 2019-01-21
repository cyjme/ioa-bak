package main

import (
	"encoding/json"
	"errors"
	"ioa"
	"ioa/proto"
	"net/http"
)

type Plugin struct {
	ioa.BasePlugin
}

type Data struct {
}

type Config struct {
	AllowOrigin      string `json:"allowOrigin"`
	AllowMethods     string `json:"allowMethods"`
	AllowHeaders     string `json:"allowHeaders"`
	ExposeHeaders    string `json:"exposeHeaders"`
	AllowCredentials string `json:"allowCredentials"`
	MaxAge           string `json:"maxAge"`
}

type RawConfig struct {
	AllowOrigin      string `json:"allowOrigin"`
	AllowMethods     string `json:"allowMethods"`
	AllowHeaders     string `json:"allowHeaders"`
	ExposeHeaders    string `json:"exposeHeaders"`
	AllowCredentials string `json:"allowCredentials"`
	MaxAge           string `json:"maxAge"`
}

func (c *Config) UnmarshalJSON(b []byte) error {
	rawConfig := RawConfig{}
	err := json.Unmarshal(b, &rawConfig)
	if err != nil {
		panic(err)
	}

	if rawConfig.AllowOrigin == "" {
		rawConfig.AllowOrigin = "*"
	}
	if rawConfig.AllowMethods == "" {
		rawConfig.AllowMethods = "POST, GET, OPTIONS, PUT, DELETE, UPDATE, PATCH"
	}
	if rawConfig.AllowHeaders == "" {
		rawConfig.AllowHeaders = "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization"
	}
	if rawConfig.ExposeHeaders == "" {
		rawConfig.ExposeHeaders = "Content-Length"
	}
	if rawConfig.AllowCredentials == "" {
		rawConfig.AllowCredentials = "true"
	}
	if rawConfig.MaxAge == "" {
		rawConfig.MaxAge = "86400"
	}

	c.AllowOrigin = rawConfig.AllowOrigin
	c.AllowMethods = rawConfig.AllowMethods
	c.AllowHeaders = rawConfig.AllowHeaders
	c.ExposeHeaders = rawConfig.ExposeHeaders
	c.AllowCredentials = rawConfig.AllowCredentials
	c.MaxAge = rawConfig.MaxAge

	return nil
}

var name = "cors"

func (i Plugin) GetName() string {
	return name
}

func (i Plugin) GetDescribe() string {
	return "set CORS"
}

func (i Plugin) GetConfigTemplate() proto.ConfigTpl {
	configTpl := proto.ConfigTpl{
		{Name: "allowOrigin", Desc: "* or http://www.test.com", Required: false, FieldType: "string"},
		{Name: "allowMethods", Desc: "POST, GET, OPTIONS, PUT, DELETE, UPDATE, PATCH", Required: false, FieldType: "string"},
		{Name: "allowHeaders", Desc: "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization", Required: false, FieldType: "string"},
		{Name: "exposeHeaders", Desc: "Content-Length", Required: false, FieldType: "string"},
		{Name: "allowCredentials", Desc: "true", Required: false, FieldType: "string"},
		{Name: "maxAge", Desc: "86400", Required: false, FieldType: "string"},
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
	i.Logger().Debug("this is config***********", config)
	api.PluginConfig[name] = config
	return nil
}

func (i Plugin) Run(w http.ResponseWriter, r *http.Request, api *ioa.Api) error {
	config := api.PluginConfig[name].(Config)

	w.Header().Set("Access-Control-Allow-Origin", config.AllowOrigin)
	w.Header().Set("Access-Control-Max-Age", config.MaxAge)
	w.Header().Set("Access-Control-Allow-Methods", config.AllowMethods)
	w.Header().Set("Access-Control-Allow-Headers", config.AllowHeaders)
	w.Header().Set("Access-Control-Expose-Headers", config.ExposeHeaders)
	w.Header().Set("Access-Control-Allow-Credentials", config.AllowCredentials)

	return nil
}

func (i Plugin) throwErr(err error) error {
	return errors.New("plugin" + name + err.Error())
}

var ExportPlugin Plugin
