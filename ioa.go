package ioa

import (
	"encoding/json"
	"ioa/httpServer/app"
	"ioa/httpServer/model"
	"ioa/router"
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"strings"
)

type Ioa struct {
	Apis    map[string]Api
	Plugins Plugins
	Router  router.Router
}

func New() *Ioa {
	return &Ioa{
		Apis:    make(map[string]Api),
		Plugins: make(Plugins),
		Router:  router.New(),
	}
}

func (ioa *Ioa) StartServer() {
	//todo 获取状态为 启用 的 api 列表，进行注册
	for _, plugin := range app.Config.Plugins {
		ioa.Plugins.Register(plugin.Name, plugin.Path)
	}

	ioa.LoadApi()
	//todo 为所有的 api 初始化插件数据

	//log.Println("为所有api 初始化插件数据")
	for _, api := range ioa.Apis {
		for _, plugin := range api.Plugins {
			log.Println("this api's plugins", api.Plugins)
			log.Println("plugin range...........", plugin)
			err := ioa.Plugins[plugin].InitApi(&api)
			if err != nil {
				log.Println(err)
			}
		}
	}

	ioa.loadApiToRouter()
	http.HandleFunc("/", ioa.ReverseProxy)

	addr := app.Config.Ioa.Host + ":" + app.Config.Ioa.Port
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		panic(err)
	}
}

func (ioa *Ioa) ReverseProxy(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	path := r.URL.Path
	apiId, params, _ := ioa.Router.FindRoute(method, path)
	log.Println("api id is ..... :", apiId)
	log.Println("api params is :", params)
	//todo match api, get model.Api
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
	log.Println("receive request")

	//todo find upstream info, and reverseProxy
	targetsLen := len(api.Targets)
	if targetsLen == 0 {
		w.Write([]byte("no target"))
		return
	}

	target := api.Targets[rand.Intn(len(api.Targets))]

	director := func(req *http.Request) {
		req.URL.Host = target.Host + ":" + target.Port
		req.URL.Scheme = target.Scheme
		req.URL.Path = target.Path
	}

	proxy := &httputil.ReverseProxy{Director: director}
	proxy.ServeHTTP(w, r)
	//todo plugin
}

func (ioa *Ioa) loadApiToRouter() {
	for id, api := range ioa.Apis {
		ioa.Router.AddRoute(strings.ToUpper(api.Method), api.Path, id)
	}
}

func (ioa *Ioa) LoadApi() {
	apiGroup := model.ApiGroup{}
	apiGroups, _, err := apiGroup.List("", "", -1, -1)
	if err != nil {
		log.Panic("list apiGroup error", err.Error())
	}

	for _, group := range *apiGroups {
		var apiGroupPolicies []string
		var apiGroupPlugins []string
		json.Unmarshal([]byte(group.Policies), &apiGroupPolicies)
		json.Unmarshal([]byte(group.Plugins), &apiGroupPlugins)
		//group policy plugins
		var groupPoliciesPlugins []string

		//解析 groupAolicies, 从中取出 plugins
		for _, policyId := range apiGroupPolicies {
			policy := model.Policy{}
			policy.Id = policyId
			policy.Get()

			var plugins []string
			json.Unmarshal([]byte(policy.Plugins), &plugins)

			groupPoliciesPlugins = append(groupPoliciesPlugins, plugins...)
		}

		for _, api := range group.Apis {
			var newApiPolicies []string
			var newApiPlugins []string
			newApiPluginsConfig := make(map[string]json.RawMessage)
			type rawPlugin struct {
				Name   string
				Config json.RawMessage
			}
			var newRawPlugins []rawPlugin
			json.Unmarshal([]byte(api.Policies), &newApiPolicies)
			json.Unmarshal([]byte(api.Plugins), &newRawPlugins)

			for _, rawPlugin := range newRawPlugins {
				newApiPlugins = append(newApiPlugins, rawPlugin.Name)
				newApiPluginsConfig[rawPlugin.Name] = rawPlugin.Config
			}

			//api policy plugins
			var apiPoliciesPlugins []string
			//解析 apiPolicies, 从中取出 plugins
			for _, policyId := range newApiPolicies {
				policy := model.Policy{}
				policy.Id = policyId
				policy.Get()
				var plugins []string
				json.Unmarshal([]byte(policy.Plugins), plugins)
				apiPoliciesPlugins = append(apiPoliciesPlugins, plugins...)
			}

			newApiAllPlugins := append(groupPoliciesPlugins, apiGroupPlugins...)
			newApiAllPlugins = append(newApiAllPlugins, apiPoliciesPlugins...)
			newApiAllPlugins = append(newApiAllPlugins, newApiPlugins...)

			newApi := Api{
				ApiGroupId: api.ApiGroupId,
				Name:       api.Name,
				Describe:   api.Describe,
				Path:       api.Path,
				Method:     api.Method,
				Status:     api.Status,

				Targets:         api.Targets,
				Params:          api.Params,
				Policies:        newApiPolicies,
				GroupPolicies:   apiGroupPolicies,
				GroupPlugins:    apiGroupPlugins,
				Plugins:         newApiPlugins,
				AllPlugin:       newApiAllPlugins,
				PluginRawConfig: newApiPluginsConfig,
				PluginConfig:    make(map[string]interface{}),
				PluginsData:     make(map[string]interface{}),
			}

			ioa.Apis[api.Id] = newApi
		}
	}
}
