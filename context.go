package ioa

import (
	"net/http"
)

type Context struct {
	Plugins        []string
	Api            Api
	Request        http.Request
	ResponseWriter http.ResponseWriter
}
