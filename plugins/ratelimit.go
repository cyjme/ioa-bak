package main

import (
	"golang.org/x/time/rate"
	"ioa/plugin"
	"log"
	"net/http"
)

type ioaPlugin struct {
}

func (s ioaPlugin) Init() {
	var limiter = rate.NewLimiter(2, 5)
}

func (s ioaPlugin) GetName() string {
	return "request_size"
}

func (s ioaPlugin) GetConfigTemplate() plugin.Config {
	config := plugin.Config{
		{Name: "maxSize", Required: true, FieldType: "int64"},
	}
	return config
}

func (s ioaPlugin) Run(w http.ResponseWriter, r *http.Request, config map[string]interface{}) {

	contentLength := r.ContentLength
	maxSize := config["maxSize"].(int64)

	if contentLength > maxSize {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("contentLength too large"))
	}

	log.Println("request content length:", contentLength)
}

var IoaPlugin ioaPlugin
