package main

import (
	"errors"
	"ioa"
	"log"
	"net/http"
	"strconv"
)

type ioaPlugin struct {
	BasePlugin
}

type Data struct {
}
type Config struct {
	maxSize int64
}

var name = "request_size"

func (s ioaPlugin) GetName() string {
	return "request_size"
}

func (s ioaPlugin) GetDescribe() string {
	return "request_size just get a request content-length"
}

func (s ioaPlugin) GetConfigTemplate() ioa.ConfigTpl {
	configTpl := ioa.ConfigTpl{
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
	maxSizeStr, exist := api.PluginRawConfig["request_size_max_size"]
	if !exist {
		return i.throwErr(errors.New("config field doesn't exist"))
	}

	maxSize, err := strconv.ParseInt(maxSizeStr, 10, 64)
	if err != nil {
		return i.throwErr(err)
	}
	config := Config{
		maxSize: maxSize,
	}
	api.PluginConfig[name] = config

	return nil
}

func (s ioaPlugin) Run(w http.ResponseWriter, r *http.Request, api *ioa.Api) error {
	contentLength := r.ContentLength
	config := api.PluginConfig[name].(Config)

	if contentLength > config.maxSize {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("contentLength too large"))
	}
	log.Println("request content length:", contentLength)

	return nil
}

func (i ioaPlugin) throwErr(err error) error {
	return errors.New("plugin" + name + err.Error())
}

var IoaPlugin ioaPlugin
