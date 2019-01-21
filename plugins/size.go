package main

import (
	"encoding/json"
	"errors"
	"ioa"
	"ioa/proto"
	"net/http"
	"strconv"
)

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

func (i Plugin) GetName() string {
	return "request_size"
}

func (i Plugin) GetDescribe() string {
	return "request_size just get a request content-length"
}

func (i Plugin) GetConfigTemplate() proto.ConfigTpl {
	configTpl := proto.ConfigTpl{
		{Name: "maxSize", Desc: "maxSize", Required: true, FieldType: "int64"},
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
	err := json.Unmarshal(api.PluginRawConfig[name], &config)
	if err != nil {
		return err
	}
	i.Logger().Debug("plugin init api config success:" + name)

	api.PluginConfig[name] = config

	return nil
}

func (i Plugin) Run(w http.ResponseWriter, r *http.Request, api *ioa.Api) error {
	contentLength := r.ContentLength
	config := api.PluginConfig[name].(Config)

	if contentLength > config.MaxSize {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("contentLength too large"))
		return errors.New("contentLength too large")
	}

	return nil
}

func (i Plugin) throwErr(err error) error {
	return errors.New("plugin" + name + err.Error())
}

var ExportPlugin Plugin
