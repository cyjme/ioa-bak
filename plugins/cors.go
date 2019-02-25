package main

import (
	"encoding/json"
	"ioa"
	"ioa/proto"
	"net/http"
)

var (
	name = "cors"
	desc = "set cors"
	tags = []string{"security"}
)

var configTpl = proto.ConfigTpl{
	{Name: "allowOrigin", Desc: "* or http://www.test.com", Required: false, FieldType: "string"},
	{Name: "allowMethods", Desc: "POST, GET, OPTIONS, PUT, DELETE, UPDATE, PATCH", Required: false, FieldType: "string"},
	{Name: "allowHeaders", Desc: "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization", Required: false, FieldType: "string"},
	{Name: "exposeHeaders", Desc: "Content-Length", Required: false, FieldType: "string"},
	{Name: "allowCredentials", Desc: "true", Required: false, FieldType: "string"},
	{Name: "maxAge", Desc: "86400", Required: false, FieldType: "string"},
}

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
		return err
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

func (i Plugin) GetName() string {
	return name
}

func (i Plugin) GetTags() []string {
	return tags
}

func (i Plugin) GetDescribe() string {
	return desc
}

func (i Plugin) GetConfigTemplate() proto.ConfigTpl {
	return configTpl
}

func (i Plugin) InitApi(api *ioa.Api) error {
	err := i.InitApiConfig(api)
	if err != nil {
		return err
	}
	err = i.InitApiData(api)
	if err != nil {
		return err
	}

	return nil
}

func (i Plugin) InitApiData(api *ioa.Api) error {
	return nil
}

func (i Plugin) InitApiConfig(api *ioa.Api) error {
	var config Config
	err := json.Unmarshal(api.PluginRawConfig[name], &config)
	if err != nil {
		return err
	}
	api.PluginConfig[name] = config

	return nil
}

func (i Plugin) ReceiveRequest(ctx *ioa.Context) {
	config := ctx.Api.PluginConfig[name].(Config)

	ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", config.AllowOrigin)
	ctx.ResponseWriter.Header().Set("Access-Control-Max-Age", config.MaxAge)
	ctx.ResponseWriter.Header().Set("Access-Control-Allow-Methods", config.AllowMethods)
	ctx.ResponseWriter.Header().Set("Access-Control-Allow-Headers", config.AllowHeaders)
	ctx.ResponseWriter.Header().Set("Access-Control-Expose-Headers", config.ExposeHeaders)
	ctx.ResponseWriter.Header().Set("Access-Control-Allow-Credentials", config.AllowCredentials)

	if ctx.Request.Method == http.MethodOptions {
		ctx.ResponseWriter.WriteHeader(http.StatusOK)
		_, err := ctx.ResponseWriter.Write(nil)
		if err != nil {
			i.Logger().Info("ResponseWriter.Writer err", err)
		}
		ctx.Cancel()
	}
}

func (i Plugin) ReceiveResponse(ctx *ioa.Context) {
	return
}

var ExportPlugin Plugin
