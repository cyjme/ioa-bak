package main

import (
	"encoding/json"
	"errors"
	"ioa"
	"ioa/proto"
	"log"
	"net/http"
	"strconv"
)

type ioaPlugin struct {
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
		panic(err)
	}

	maxSize, err := strconv.ParseInt(rawConfig.MaxSize, 10, 64)
	if err != nil {
		panic(err)
	}

	c.MaxSize = maxSize
	return nil
}

var name = "request_size"

func (s ioaPlugin) GetName() string {
	return "request_size"
}

func (s ioaPlugin) GetDescribe() string {
	return "request_size just get a request content-length"
}

func (s ioaPlugin) GetConfigTemplate() proto.ConfigTpl {
	configTpl := proto.ConfigTpl{
		{Name: "maxSize", Desc: "maxSize", Required: true, FieldType: "int64"},
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
	contentLength := r.ContentLength
	config := api.PluginConfig[name].(Config)

	if contentLength > config.MaxSize {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("contentLength too large"))
		return errors.New("contentLength too large")
	}
	log.Println("request content length:", contentLength)

	return nil
}

func (i ioaPlugin) throwErr(err error) error {
	return errors.New("plugin" + name + err.Error())
}

var IoaPlugin ioaPlugin
