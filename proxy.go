package ioa

import (
	"io"
	"math/rand"
	"net/http"
	"strings"
)

func (ioa *Ioa) reverseProxy(w http.ResponseWriter, r *http.Request) {
	method := r.Method

	//todo
	if method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		w.Write(nil)
		return
	}

	path := r.URL.Path
	apiId, params, _ := ioa.Router.FindRoute(method, path)

	if apiId == "" {
		w.WriteHeader(404)
		w.Write([]byte("no find the router"))
		return
	}

	//todo match api, get store.Api
	api := ioa.Apis[apiId]

	ctx := Context{
		ResponseWriter: w,
		Request:        r,
		Response:       nil,
		Api:            &api,
		Next:           true,
	}

	for _, plugin := range api.Plugins {
		plugin, exist := ioa.Plugins[plugin]
		if !exist {
			w.Write([]byte("the api use unexist plugin:" + plugin.GetName()))
			return
		}

		plugin.ReceiveRequest(&ctx)
		if !ctx.Next {
			return
		}
	}

	targetsLen := len(api.Targets)
	if targetsLen == 0 {
		w.Write([]byte("no target"))
		return
	}

	target := api.Targets[rand.Intn(len(api.Targets))]

	targetPath := target.Path
	for _, param := range params {
		targetPath = strings.Replace(targetPath, ":"+param.Key, param.Value, 1)
		targetPath = strings.Replace(targetPath, "*"+param.Key, param.Value, 1)
	}

	url := target.Scheme + target.Host + ":" + target.Port + targetPath
	if target.Method == "*" || target.Method == "" {
		target.Method = r.Method
	}
	newReq, err := http.NewRequest(strings.ToUpper(target.Method), url, r.Body)
	newReq.Header = r.Header
	if err != nil {
		log.Debug(ERR_PROXY_CREATE_REQUEST, err)
	}

	resp, err := client.Do(newReq)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		log.Debug(ERR_PROXY_DO_REQUEST, err)
	}

	ctx.Response = resp
	for _, plugin := range api.Plugins {
		plugin, exist := ioa.Plugins[plugin]
		if !exist {
			w.Write([]byte("the api use unexist plugin :" + plugin.GetName()))
			return
		}
		plugin.ReceiveResponse(&ctx)
		if !ctx.Next {
			return
		}
	}

	for name, values := range resp.Header {
		w.Header()[name] = values
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
	return
}
