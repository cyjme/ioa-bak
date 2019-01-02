package main

import (
	"errors"
	"golang.org/x/time/rate"
	"ioa"
	"ioa/httpServer/app"
	"log"
	"net/http"
	"strconv"
)

type ioaPlugin struct {
}

type Data struct {
	Limiter *rate.Limiter
}
type Config struct {
	limit rate.Limit
	burst int
}

var name = "rate_limit"

func (i ioaPlugin) GetName() string {
	return name
}

func (i ioaPlugin) GetConfigTemplate() ioa.ConfigTpl {
	configTpl := ioa.ConfigTpl{
		{Name: "Limit", Desc: "The number of events per second.", Required: true, FieldType: "float64"},
	}

	return configTpl
}

func (i ioaPlugin) Run(w http.ResponseWriter, r *http.Request, api *ioa.Api) error {
	//limit := config["limit"].(float64)
	log.Println("rate limiter plugin run")
	log.Println("now count is ", app.Count)

	data := api.PluginsData[name].(Data)
	if !data.Limiter.Allow() {
		w.Write([]byte("rate limit"))
		log.Println("not allow")
		return errors.New("rate limit")
	}

	return nil
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
	config := api.PluginConfig[name].(Config)
	var Limiter = rate.NewLimiter(config.limit, config.burst)

	api.PluginsData[name] = Data{
		Limiter: Limiter,
	}

	return nil
}

func (i ioaPlugin) InitApiConfig(api *ioa.Api) error {
	limitString, exist := api.PluginRawConfig["rate_limit_limit"]
	if !exist {
		return i.throwErr(errors.New("config field doesn't exist"))
	}
	burstString, exist := api.PluginRawConfig["rate_limit_burst"]

	if !exist {
		return i.throwErr(errors.New("config field doesn't exist"))
	}
	burst, err := strconv.Atoi(burstString)
	if err != nil {
		return i.throwErr(err)
	}

	limit, err := strconv.Atoi(limitString)
	if err != nil {
		return i.throwErr(err)
	}

	config := Config{
		limit: rate.Limit(limit),
		burst: burst,
	}
	api.PluginConfig[name] = config

	return nil
}

func (i ioaPlugin) throwErr(err error) error {
	return errors.New("plugin" + name + err.Error())
}

var IoaPlugin ioaPlugin
