package main

import (
	"errors"
	"ioa/httpServer/app"
	"ioa/plugin"
	"log"
	"net/http"
)

type ioaPlugin struct {
}

func (i ioaPlugin) GetName() string {
	return "rate_limit"
}

func (i ioaPlugin) GetConfigTemplate() plugin.ConfigTpl {
	configTpl := plugin.ConfigTpl{
		{Name: "Limit", Desc: "The number of events per second.", Required: true, FieldType: "float64"},
	}

	return configTpl
}

func (i ioaPlugin) Run(w http.ResponseWriter, r *http.Request, config map[string]interface{}) error {
	//limit := config["limit"].(float64)
	log.Println("rate limiter plugin run")
	app.Count++
	log.Println("now count is ", app.Count)

	if !app.Limiter.Allow() {
		w.Write([]byte("rate limit"))
		log.Println("not allow")
		return errors.New("rate limit")
	}

	return nil
}

var IoaPlugin ioaPlugin
