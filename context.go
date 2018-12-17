package main

import (
	"ioa/httpServer/model"
	"net/http"
)

type Context struct {
	Api            model.Api
	Request        http.Request
	ResponseWriter http.ResponseWriter
}
