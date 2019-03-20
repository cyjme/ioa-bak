package main

import (
	"encoding/json"
	"ioa"
	"ioa/proto"
	"net/http"
	"strconv"
	"strings"
)

var (
	name = "default_response"
	desc = "default_response, if use the plugin, make sure it is the first response plugin"
	tags = []string{"response"}
)

type Plugin struct {
	ioa.BasePlugin
}

type Data struct {
}

type Config struct {
	StatusCode int               `json:"statusCode"`
	Headers    map[string]string `json:"headers"`
	Content    string            `json:"content"`
}

type RawConfig struct {
	StatusCode string `json:"statusCode"`
	Headers    string `json:"headers"`
	Content    string `json:"content"`
}

func (c *Config) UnmarshalJSON(b []byte) error {
	rawConfig := RawConfig{}
	err := json.Unmarshal(b, &rawConfig)
	if err != nil {
		panic(err)
	}
	c.Headers = make(map[string]string)

	headerValues := strings.Split(rawConfig.Headers, ",")
	for _, headerValue := range headerValues {
		headerValueArray := strings.Split(headerValue, ":")
		header := headerValueArray[0]
		value := headerValueArray[1]
		c.Headers[header] = value
	}

	c.StatusCode, err = strconv.Atoi(rawConfig.StatusCode)
	if err != nil {
		return err
	}
	c.Content = rawConfig.Content

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
	configTpl := proto.ConfigTpl{
		{Name: "statusCode", Desc: "status code", Required: true, FieldType: "int"},
		{Name: "headers", Desc: "response header array, separated by , e.g.: content-type:application/json,age:3600", Required: true, FieldType: "string"},
		{Name: "content", Desc: "response content", Required: true, FieldType: "string"},
	}

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
	return
}

func (i Plugin) ReceiveResponse(ctx *ioa.Context) {
	if ctx.Response.StatusCode == 200 {
		return
	}
	if ctx.Response.Header == nil {
		ctx.Response.Header = make(http.Header)
	}

	config := ctx.Api.PluginConfig[name].(Config)
	for headerName, headerValue := range config.Headers {
		ctx.ResponseWriter.Header().Set(headerName, headerValue)
	}
	ctx.ResponseWriter.WriteHeader(config.StatusCode)
	ctx.ResponseWriter.Write([]byte(config.Content))
	ctx.Cancel()
}

var ExportPlugin Plugin
