package main

import (
	"ioa/httpServer"
	"ioa/plugin"
	"log"
	"net/http"
)

type Ioa struct {
}

type Api struct {

}

func main() {
	//plugin.Load()
	go httpServer.Run()
	StartServer()
}

func StartServer() {
	//todo load plugin
	pluginCenter := plugin.NewPluginCenter()
	pluginCenter.Register("001", "./plugins/count.so")

	http.HandleFunc("/", ReverseProxy)
	http.ListenAndServe(":11112", nil)
}

func ReverseProxy(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
	log.Println("receive request")
	//todo  use httpRouter get apiId

	//todo match api, get model.Api
	//todo plugin
}

//func LoadApi() {
//	apiGroup := model.ApiGroup{}
//	apiGroups, _, err := apiGroup.List("", "", -1, -1)
//	if err != nil {
//		log.Panic("list apiGroup error", err.Error())
//	}
//	for _, group := range apiGroups {
//
//	}
//
//}
