package main

import (
	"encoding/json"
	"ioa"
	"ioa/proto"
	"net/http"
	"strconv"
)

var (
	name = "request_size"
	desc = "request_size just get a request content-length"
	tags = []string{"traffic control"}

	RESP_CONTENT_TOO_LARGE = "contentLength too large"
)

var configTpl = proto.ConfigTpl{
	{Name: "maxSize", Desc: "maxSize", Required: true, FieldType: "int64"},
}

type Plugin struct {
	ioa.BasePlugin
}

type Data struct {
}

type Config struct {
	MaxSize int64 `json:"maxSize"`
}

type RawConfig struct {
	MaxSize string `json:"maxSize"`
}

func (c *Config) UnmarshalJSON(b []byte) error {
	rawConfig := RawConfig{}
	err := json.Unmarshal(b, &rawConfig)
	if err != nil {
		return err
	}

	maxSize, err := strconv.ParseInt(rawConfig.MaxSize, 10, 64)
	if err != nil {
		return err
	}

	c.MaxSize = maxSize
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
	contentLength := ctx.Request.ContentLength
	config := ctx.Api.PluginConfig[name].(Config)

	if contentLength > config.MaxSize {
		ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		ctx.ResponseWriter.Write([]byte(RESP_CONTENT_TOO_LARGE))
		ctx.Cancel()
	}
}

func (i Plugin) ReceiveResponse(ctx *ioa.Context) {
}

var ExportPlugin Plugin
