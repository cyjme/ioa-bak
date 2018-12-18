package main

import (
	"log"
	"net/http"
)

func Run(w http.ResponseWriter, r *http.Request) {
	log.Println("plugin=size")
	contentLength := r.ContentLength

	log.Println("request content length:", contentLength)
}
