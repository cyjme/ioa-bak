package main

import (
	"ioa/plugin"
	"log"
	"net/http"
)

type ioaPlugin struct {
}

func (s ioaPlugin) GetName() string {
	return "request_size"
}

func (s ioaPlugin) GetConfigTemplate() plugin.ConfigTpl {
	configTpl := plugin.ConfigTpl{
		{Name: "maxSize", Desc: "maxSize", Required: true, FieldType: "int64"},
	}

	return configTpl
}

func (s ioaPlugin) Run(w http.ResponseWriter, r *http.Request, config map[string]interface{}) error{
	contentLength := r.ContentLength
	maxSize := config["maxSize"].(int64)

	if contentLength > maxSize {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("contentLength too large"))
	}

	log.Println("request content length:", contentLength)

	return nil
}

var IoaPlugin ioaPlugin
