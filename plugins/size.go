package main

import (
	"log"
	"net/http"
)

type sizePlugin struct {
}

func (s sizePlugin) GetName() string {
	return "request_size"
}

func (s sizePlugin) GetConfig() {

}

func (s sizePlugin) Run(w http.ResponseWriter, r *http.Request) {
	contentLength := r.ContentLength

	log.Println("request content length:", contentLength)
}

var SizePlugin sizePlugin
