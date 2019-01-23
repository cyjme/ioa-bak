package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"ioa"
	"ioa/proto"
	"net/http"
	"strings"
)

var (
	name = "copy_request"
	desc = "copy_request to new url"
	tags = []string{"traffic control"}
)

var configTpl = proto.ConfigTpl{
	{Name: "urls", Desc: "urls split by ,", Required: true, FieldType: "string"},
}

type Plugin struct {
	ioa.BasePlugin
}

type Data struct {
}
type Config struct {
	Urls []string `json:"urls"`
}

type RawConfig struct {
	Urls string `json:"urls"`
}

func (c *Config) UnmarshalJSON(b []byte) error {
	rawConfig := RawConfig{}
	err := json.Unmarshal(b, &rawConfig)
	if err != nil {
		panic(err)
	}
	c.Urls = strings.Split(rawConfig.Urls, ",")

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

	for _, url := range config.Urls {
		err := doRequest(ctx.Request, url)
		if err != nil {
			return
		}
	}

	return
}

func doRequest(r *http.Request, url string) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	newReq, err := http.NewRequest(r.Method, url, bytes.NewReader(body))
	for header, value := range r.Header {
		newReq.Header[header] = value
	}
	client := http.Client{}

	_, err = client.Do(newReq)
	if err != nil {
		return err
	}

	return nil
}

func (i Plugin) ReceiveResponse(ctx *ioa.Context) {
	return
}

var ExportPlugin Plugin
