package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"ioa"
	"ioa/proto"
	"net/http"
	"strings"
)

type ioaPlugin struct {
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

const name = "copy_request"

func (i ioaPlugin) GetName() string {
	return name
}

func (i ioaPlugin) GetDescribe() string {
	return "copy_request to new url"
}

func (i ioaPlugin) GetConfigTemplate() proto.ConfigTpl {
	configTpl := proto.ConfigTpl{
		{Name: "urls", Desc: "urls split by ,", Required: true, FieldType: "string"},
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
	err := json.Unmarshal(api.PluginRawConfig[name], &config)
	if err != nil {
		return err
	}
	i.Logger().Debug("plugin init api config success,plugin: " + name + "api: " + api.Name)

	api.PluginConfig[name] = config

	return nil
}

func (i ioaPlugin) Run(w http.ResponseWriter, r *http.Request, api *ioa.Api) error {
	config := api.PluginConfig[name].(Config)

	for _, url := range config.Urls {
		err := doRequest(r, url)
		if err != nil {
			return err
		}
	}

	return nil
}

func (i ioaPlugin) throwErr(err error) error {
	return errors.New("plugin" + name + err.Error())
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

var IoaPlugin ioaPlugin
