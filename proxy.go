package ioa

import (
	"io"
	"math/rand"
	"net/http"
	"strings"
)

func (ioa *Ioa) reverseProxy(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	path := r.URL.Path
	apiId, _, _ := ioa.Router.FindRoute(method, path)

	if apiId == "" {
		w.WriteHeader(404)
		w.Write([]byte("no find the router"))
		return
	}

	//todo match api, get store.Api
	api := ioa.Apis[apiId]

	for _, plugin := range api.Plugins {
		plugin, exist := ioa.Plugins[plugin]
		if !exist {
			w.Write([]byte("the api use unexist plugin:" + plugin.GetName()))
			return
		}

		name := plugin.GetName()
		log.Debug("plugin will run", name)
		err := plugin.Run(w, r, &api)
		if err != nil {
			log.Error("plugin run error :", err)
			w.WriteHeader(http.StatusBadGateway)
			w.Write([]byte("gateway error plugin run error"))
			return
		}
	}

	//todo find upstream info, and reverseProxy
	targetsLen := len(api.Targets)
	if targetsLen == 0 {
		w.Write([]byte("no target"))
		return
	}

	target := api.Targets[rand.Intn(len(api.Targets))]

	url := target.Scheme + target.Host + ":" + target.Port + target.Path
	newReq, err := http.NewRequest(strings.ToUpper(target.Method), url, r.Body)
	newReq.Header = r.Header
	if err != nil {
		log.Debug("err", err)
	}

	resp, err := client.Do(newReq)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		log.Debug("err", err)
	}

	for name, values := range resp.Header {
		w.Header()[name] = values
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
	return
}
