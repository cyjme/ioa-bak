package ioa

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"ioa/proto"
	"ioa/router"
	"ioa/store"
	"log"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"strings"
)

type Ioa struct {
	Apis    map[string]Api
	Plugins Plugins
	Router  router.Router
	Config  Config
}

var client *http.Client

func New(config Config) *Ioa {
	return &Ioa{
		Apis:    make(map[string]Api),
		Plugins: make(Plugins),
		Router:  router.New(),
		Config:  config,
	}
}

func (ioa *Ioa) StartServer() {
	// Customize the Transport to have larger connection pool
	defaultRoundTripper := http.DefaultTransport
	defaultTransportPointer, ok := defaultRoundTripper.(*http.Transport)
	if !ok {
		panic(fmt.Sprintf("defaultRoundTripper not an *http.Transport"))
	}
	defaultTransport := *defaultTransportPointer // dereference it to get a copy of the struct that the pointer points to
	defaultTransport.MaxIdleConns = ioa.Config.Proxy.MaxIdleConns
	defaultTransport.MaxIdleConnsPerHost = ioa.Config.Proxy.MaxIdleConnsPerHost
	client = &http.Client{Transport: &defaultTransport}

	http.HandleFunc("/", ioa.ReverseProxy)
	ioa.Load()
	go ioa.Watch()

	addr := ioa.Config.Proxy.Host + ":" + ioa.Config.Proxy.Port
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		panic(err)
	}
}

func (ioa *Ioa) Watch() {
	api := store.Api{}
	go api.Watch(func() {
		ioa.Load()
	})
}

func (ioa *Ioa) Load() {
	log.Println("ioa start Load data")
	ioa.Apis = make(map[string]Api)
	ioa.Plugins = make(Plugins)
	ioa.Router.Clear()

	//获取状 api 列表,注册到 ioa
	for _, plugin := range ioa.Config.Plugins {
		ioa.Plugins.Register(plugin.Name, plugin.Path)
	}

	//注册插件信息到 store
	var plugins []proto.Plugin
	for _, plugin := range ioa.Plugins {
		plugins = append(plugins, proto.Plugin{
			Name:      plugin.GetName(),
			Describe:  plugin.GetDescribe(),
			ConfigTpl: plugin.GetConfigTemplate(),
		})
	}
	log.Println("ReCreate plugins info to store")
	store.ReCreatePlugin(plugins)

	//从数据库中加载 api
	ioa.LoadApi()

	//为所有的 api 使用的插件,初始化 config 和 api 缓存数据
	for _, api := range ioa.Apis {
		for _, plugin := range api.Plugins {
			err := ioa.Plugins[plugin].InitApi(&api)
			if err != nil {
				log.Println(err)
			}
		}
	}

	//把 api 加载到router中
	ioa.loadApiToRouter()
}

func (ioa *Ioa) ReverseProxy(w http.ResponseWriter, r *http.Request) {
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
		log.Println("plugin will run", name)
		err := plugin.Run(w, r, &api)
		if err != nil {
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
	newReq, err := http.NewRequest(target.Method, url, r.Body)
	if err != nil {
		log.Println("err", err)
	}

	resp, err := client.Do(newReq)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		log.Println("err", err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("err", err)
	}

	w.Header().Set("content-type", resp.Header.Get("content-type"))
	w.Write(body)
	return

	//使用reverseProxy 导致连接不能复用。
	//director := func(req *http.Request) {
	//	req.URL.Host = target.Host + ":" + target.Port
	//	req.URL.Scheme = target.Scheme
	//	req.URL.Path = target.Path
	//}
	//
	//proxy := &httputil.ReverseProxy{Director: director}
	//proxy.ServeHTTP(w, r)

	//todo plugin
}

func transferRequest(oldReq *http.Request) {
}

func (ioa *Ioa) loadApiToRouter() {
	for id, api := range ioa.Apis {
		ioa.Router.AddRoute(strings.ToUpper(api.Method), api.Path, id)
	}
}

func (ioa *Ioa) LoadApi() {
	api := store.Api{}
	apis, _, err := api.List()
	if err != nil {
		panic(err)
	}

	var newApiPlugins []string
	newApiPluginsConfig := make(map[string]json.RawMessage)

	type rawPlugin struct {
		Name   string
		Config json.RawMessage
	}
	for _, api := range apis {
		var newRawPlugins []rawPlugin

		if api.Plugins == "" {
			api.Plugins = "[]"
		}

		if err := json.Unmarshal([]byte(api.Plugins), &newRawPlugins); err != nil {
			log.Println(err)
			continue
		}

		for _, rawPlugin := range newRawPlugins {
			newApiPlugins = append(newApiPlugins, rawPlugin.Name)
			newApiPluginsConfig[rawPlugin.Name] = rawPlugin.Config
		}

		newApi := Api{
			ApiGroupId: api.ApiGroupId,
			Name:       api.Name,
			Describe:   api.Describe,
			Path:       api.Path,
			Method:     api.Method,
			Status:     api.Status,

			Targets:         api.Targets,
			Plugins:         newApiPlugins,
			PluginRawConfig: newApiPluginsConfig,
			PluginConfig:    make(map[string]interface{}),
			PluginsData:     make(map[string]interface{}),
		}

		ioa.Apis[api.Id] = newApi
	}
}
