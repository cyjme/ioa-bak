package main

import (
	"encoding/json"
	"io/ioutil"
	"ioa"
	"ioa/proto"
	"net/http"
	"strconv"
	"time"
)

var (
	name = "cache"
	desc = "cache response, request path is identification"
	tags = []string{"traffic control"}
)
var configTpl = proto.ConfigTpl{
	{Name: "Enable", Desc: "enable api cache", Required: true, FieldType: "bool"},
	{Name: "TTL", Desc: "response data time to live, in seconds", Required: true, FieldType: "uint"},
}

type Plugin struct {
	ioa.BasePlugin
}

type Cache struct {
	Response       Response
	LastUpdateTime time.Time
}

type Response struct {
	Header     http.Header
	StatusCode int
	Body       []byte
}
type Data map[string]Cache

type Config struct {
	Enable bool          `json:"enable"`
	TTL    time.Duration `json:"ttl"`
}

type RawConfig struct {
	Enable string `json:"enable"`
	TTL    string `json:"ttl"`
}

func (c *Config) UnmarshalJSON(b []byte) error {
	rawConfig := RawConfig{}
	err := json.Unmarshal(b, &rawConfig)
	if err != nil {
		return err
	}
	c.Enable = false
	if rawConfig.Enable == "1" {
		c.Enable = true
	}

	ttlInt, err := strconv.Atoi(rawConfig.TTL)
	if err != nil {
		return err
	}
	c.TTL = time.Second * time.Duration(uint64(ttlInt))

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
	api.PluginsData[name] = make(Data)
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

	if !config.Enable {
		return
	}

	data := ctx.Api.PluginsData[name].(Data)
	cache, exist := data[ctx.Request.URL.Path]
	if !exist {
		return
	}
	now := time.Now()
	if now.Sub(cache.LastUpdateTime) < config.TTL {
		for name, values := range cache.Response.Header {
			ctx.ResponseWriter.Header()[name] = values
		}
		ctx.ResponseWriter.WriteHeader(cache.Response.StatusCode)
		_, err := ctx.ResponseWriter.Write(cache.Response.Body)
		if err != nil {
		}
		ctx.Cancel()
	}
}

func (i Plugin) ReceiveResponse(ctx *ioa.Context) {
	data := ctx.Api.PluginsData[name].(Data)
	body, err := ioutil.ReadAll(ctx.Response.Body)
	if err != nil {
		ctx.Cancel()
		return
	}
	data[ctx.Request.URL.Path] = Cache{
		Response: Response{
			Header:     ctx.Response.Header,
			Body:       body,
			StatusCode: ctx.Response.StatusCode,
		},
		LastUpdateTime: time.Now(),
	}
	ctx.Api.PluginsData[name] = data
}

var ExportPlugin Plugin
