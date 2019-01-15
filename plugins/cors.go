package main

import (
	"encoding/json"
	"errors"
	"ioa"
	"ioa/proto"
	"log"
	"net/http"
)

type ioaPlugin struct {
}

type Data struct {
}

type Config struct {
	AllowOrigin      string `json:"allowOrigin"`      //"*"
	AllowMethods     string `json:"allowMethods"`     //"POST, GET, OPTIONS, PUT, DELETE, UPDATE, PATCH"
	AllowHeaders     string `json:"allowHeaders"`     //"Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Auth"
	ExposeHeaders    string `json:"exposeHeaders"`    //"Content-Length"
	AllowCredentials string `json:"allowCredentials"` //"true"
	MaxAge           string `json:"maxAge"`           //"86400"
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

func (s ioaPlugin) GetName() string {
	return name
}

func (s ioaPlugin) GetDescribe() string {
	return "set CORS"
}

func (s ioaPlugin) GetConfigTemplate() proto.ConfigTpl {
	configTpl := proto.ConfigTpl{
		{Name: "allowOrigin", Desc: "", Required: false, FieldType: "string"},
		{Name: "allowMethods", Desc: "", Required: false, FieldType: "string"},
		{Name: "allowHeaders", Desc: "", Required: false, FieldType: "string"},
		{Name: "exposeHeaders", Desc: "", Required: false, FieldType: "string"},
		{Name: "allowCredentials", Desc: "", Required: false, FieldType: "string"},
		{Name: "maxAge", Desc: "", Required: false, FieldType: "string"},
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
	config := api.PluginConfig[name].(Config)

	w.Header().Set("Access-Control-Allow-Origin", config.AllowOrigin)
	w.Header().Set("Access-Control-Max-Age", config.MaxAge)
	w.Header().Set("Access-Control-Allow-Methods", config.AllowMethods)
	w.Header().Set("Access-Control-Allow-Headers", config.AllowHeaders)
	w.Header().Set("Access-Control-Expose-Headers", config.ExposeHeaders)
	w.Header().Set("Access-Control-Allow-Credentials", config.AllowCredentials)

	return nil
}

func (i ioaPlugin) throwErr(err error) error {
	return errors.New("plugin" + name + err.Error())
}

var IoaPlugin ioaPlugin
