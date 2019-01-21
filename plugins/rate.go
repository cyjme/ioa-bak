package main

import (
	"encoding/json"
	"errors"
	"golang.org/x/time/rate"
	"ioa"
	"ioa/proto"
	"strconv"
)

type Plugin struct {
	ioa.BasePlugin
}

type Data struct {
	Limiter *rate.Limiter
}
type Config struct {
	Limit rate.Limit `json:"limit"`
	Burst int        `json:"burst"`
}

var name = "rate_limit"

const DEFAULT_BURST = "0"
const DEFAULT_LIMIT = "0"

type RawConfig struct {
	Limit string `json:"limit"`
	Burst string `json:"burst"`
}

func (c *Config) UnmarshalJSON(b []byte) error {
	rawConfig := RawConfig{}
	err := json.Unmarshal(b, &rawConfig)
	if err != nil {
		return err
	}
	if rawConfig.Burst == "" {
		rawConfig.Burst = DEFAULT_BURST
	}
	if rawConfig.Limit == "" {
		rawConfig.Limit = DEFAULT_LIMIT
	}
	limitInt, err := strconv.Atoi(rawConfig.Limit)
	if err != nil {
		return err
	}
	burstInt, err := strconv.Atoi(rawConfig.Burst)
	if err != nil {
		return err
	}
	c.Limit = rate.Limit(limitInt)
	c.Burst = burstInt

	return nil
}

func (i Plugin) GetName() string {
	return name
}

func (i Plugin) GetDescribe() string {
	return `rate_limit describe`
}

func (i Plugin) GetConfigTemplate() proto.ConfigTpl {
	configTpl := proto.ConfigTpl{
		{Name: "limit", Desc: "The number of events per second.", Required: true, FieldType: "float64"},
		{Name: "burst", Desc: "The number of events for burst", Required: true, FieldType: "float64"},
	}

	return configTpl
}

func (i Plugin) Run(ctx ioa.Context) error {
	//limit := config["limit"].(float64)

	data := ctx.Api.PluginsData[name].(Data)
	if !data.Limiter.Allow() {
		ctx.ResponseWriter.Write([]byte("rate limit"))
		return errors.New("rate limit")
	}

	return nil
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
	config := api.PluginConfig[name].(Config)
	var Limiter = rate.NewLimiter(config.Limit, config.Burst)

	api.PluginsData[name] = Data{
		Limiter: Limiter,
	}

	return nil
}

func (i Plugin) InitApiConfig(api *ioa.Api) error {
	var config Config
	err := json.Unmarshal(api.PluginRawConfig[name], &config)

	if err != nil {
		return err
	}
	api.PluginConfig[name] = config
	i.Logger().Debug("plugin init api config success:" + name)

	return nil
}

func (i Plugin) throwErr(err error) error {
	return errors.New("plugin" + name + err.Error())
}

var ExportPlugin Plugin
