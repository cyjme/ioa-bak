package main

import (
	"ioa/httpServer/model"
	"ioa/router/httprouter"
)

var Router *httprouter.Router
var RouteMap = make(map[string]string)
var BackendRouteMap = make(map[string]model.Api)
