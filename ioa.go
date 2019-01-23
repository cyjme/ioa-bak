package ioa

import (
	"encoding/json"
	"fmt"
	logger "ioa/log"
	"ioa/monitor"
	"ioa/proto"
	"ioa/router"
	"ioa/store"
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
var log = logger.Get()

func New(config Config) *Ioa {
	return &Ioa{
		Apis:    make(map[string]Api),
		Plugins: make(Plugins),
		Router:  router.New(),
		Config:  config,
	}
}

func (ioa *Ioa) StartServer() {
	defaultRoundTripper := http.DefaultTransport
	defaultTransportPointer, ok := defaultRoundTripper.(*http.Transport)
	if !ok {
		panic(fmt.Sprintf("defaultRoundTripper not an *http.Transport"))
	}
	defaultTransport := *defaultTransportPointer // dereference it to get a copy of the struct that the pointer points to
	defaultTransport.MaxIdleConns = ioa.Config.Proxy.MaxIdleConns
	defaultTransport.MaxIdleConnsPerHost = ioa.Config.Proxy.MaxIdleConnsPerHost
	client = &http.Client{Transport: &defaultTransport}
	http.HandleFunc("/", ioa.reverseProxy)
	http.HandleFunc("/monitor", monitor.Handle)
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
	log.Info("start load apis/plugins from store")
	ioa.Apis = make(map[string]Api)
	ioa.Plugins = make(Plugins)
	ioa.Router.Clear()

	//获取状 plugins config,注册到 ioa
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
	log.Info("push plugins data to etcd store")
	store.ReCreatePlugin(plugins)

	//从数据库中加载 api
	ioa.LoadApi()

	//为所有的 api 使用的插件,初始化 config 和 api 缓存数据
	for _, api := range ioa.Apis {
		for _, plugin := range api.Plugins {
			err := ioa.Plugins[plugin].InitApi(&api)
			if err != nil {
				log.Error(ioa.Plugins[plugin].GetName(), ERR_INIT_API_PLUGIN, err)
			}
		}
	}

	//把 api 加载到router中
	ioa.loadApiToRouter()
}

func (ioa *Ioa) loadApiToRouter() {
	for id, api := range ioa.Apis {
		ioa.Router.AddRoute(strings.ToUpper(api.Method), api.Path, id)
	}
}

func (ioa *Ioa) LoadApi() {
	api := store.Api{}
	apis, _, err := api.List()
	log.Debug("read api from store: ", apis)
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
			log.Error(ERR_API_GET_PLUGINS, err)
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

	log.Debug("load api from store to ioa", ioa.Apis)
}
