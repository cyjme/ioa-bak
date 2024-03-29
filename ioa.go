package ioa

import (
	"encoding/json"
	"fmt"
	logger "ioa/log"
	"ioa/monitor"
	"ioa/proto"
	"ioa/router"
	"ioa/store"
	"ioa/util"
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
	log.Info("ioa proxy server starting, listen:" + addr)
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
	policy := store.Policy{}
	go policy.Watch(func() {
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
			Tags:      plugin.GetTags(),
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
		err := ioa.Router.AddRoute(strings.ToUpper(api.Method), api.Path, id)
		if err != nil {
			log.Error(err)
		}
	}
}

type pluginNames []string
type pluginRawConfigs map[string]json.RawMessage

func (ioa *Ioa) LoadApi() {
	api := store.Api{}
	apis, _, err := api.List()
	log.Debug("read api from store: ", apis)
	if err != nil {
		panic(err)
	}

	for _, api := range apis {
		//获取 api 中定义的 plugin 和 rawConfig
		apiPluginNames, apiPluginRawConfigs, err := convPluginsString(api.Plugins)
		if err != nil {
			log.Error(ERR_INIT_API_PLUGIN)
			continue
		}

		//获取 api 中的 policies 中定义的 plugins, 和 api 的进行合并，api 中的 plugin 优先级更高
		var apiPoliciesPluginNames pluginNames
		var apiPoliciesPluginRawConfigs pluginRawConfigs
		for _, policy := range api.PoliciesData {
			plugins := policy.Plugins
			apiPoliciesPluginNames, apiPoliciesPluginRawConfigs, err = convPluginsString(plugins)

			if err != nil {
				log.Error(ERR_INIT_API_PLUGIN)
				continue
			}

			for _, apiPoliciesPluginName := range apiPoliciesPluginNames {
				if inArray(apiPluginNames, apiPoliciesPluginName) {
					continue
				}
				apiPluginNames = append(apiPluginNames, apiPoliciesPluginName)
				apiPluginRawConfigs[apiPoliciesPluginName] = apiPoliciesPluginRawConfigs[apiPoliciesPluginName]
			}
		}

		for _, method := range api.Methods {
			newApi := Api{
				ApiGroupId:      api.ApiGroupId,
				Name:            api.Name,
				Describe:        api.Describe,
				Path:            api.Path,
				Method:          method,
				Status:          api.Status,
				Targets:         api.Targets,
				Plugins:         apiPluginNames,
				PluginRawConfig: apiPluginRawConfigs,
				PluginConfig:    make(map[string]interface{}),
				PluginsData:     make(map[string]interface{}),
			}

			ioa.Apis[util.GetRandomString(11)] = newApi
		}
	}

	log.Debug("load api from store to ioa", ioa.Apis)
}

func convPluginsString(plugins string) (pluginNames, pluginRawConfigs, error) {
	type rawPlugin struct {
		Name   string
		Config json.RawMessage
	}

	var newRawPlugins []rawPlugin
	var newApiPlugins []string
	newApiPluginsConfig := make(map[string]json.RawMessage)

	if plugins == "" {
		plugins = "[]"
	}

	if err := json.Unmarshal([]byte(plugins), &newRawPlugins); err != nil {
		log.Error(ERR_API_GET_PLUGINS, err)
		return nil, nil, err
	}

	for _, rawPlugin := range newRawPlugins {
		newApiPlugins = append(newApiPlugins, rawPlugin.Name)
		newApiPluginsConfig[rawPlugin.Name] = rawPlugin.Config
	}
	return newApiPlugins, newApiPluginsConfig, nil
}

func inArray(array []string, val string) bool {
	for _, item := range array {
		if item == val {
			return true
		}
	}

	return false
}
