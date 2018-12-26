package router

import (
	"ioa/router/httprouter"
)

type Router struct {
	router *httprouter.Router
}

func New() Router {
	return Router{
		router: httprouter.New(),
	}
}

func (r Router) AddRoute(method string, path string, routeId string) {
	r.router.Handle(method, path, routeId)
}

func (r Router) FindRoute(method string, path string) (string, httprouter.Params, bool) {
	routeId, param, tsr := r.router.Lookup(method, path)

	return routeId, param, tsr
}

func (r Router) Refresh() {
	r.router.ClearRoute("GET")
	r.router.ClearRoute("POST")
	r.router.ClearRoute("PUT")
	r.router.ClearRoute("DELETE")
	r.router.ClearRoute("PATCH")
	r.router.ClearRoute("OPTION")
	r.router.ClearRoute("HEADER")
}
