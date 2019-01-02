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
	log.Println("load Api from database")
	//todo 获取状态为 启用 的 api 列表，进行注册
	ioa.Plugins.Register("request_size", "./plugins/size.so")
	ioa.Plugins.Register("rate_limiter", "./plugins/rate.so")

	ioa.LoadApi()
	//todo 为所有的 api 初始化插件数据

	//log.Println("为所有api 初始化插件数据")
	//for _, api := range ioa.Apis {
	//	for _, plugin := range ioa.Plugins {
	//		err := plugin.InitApi(&api)
	//		if err != nil {
	//			log.Println(err)
	//		}
	//	}
	//}

	ioa.loadApiToRouter()
	log.Println("load Api from database Success:", ioa.Apis)
	http.HandleFunc("/", ioa.ReverseProxy)

	addr := app.Config.Ioa.Host + ":" + app.Config.Ioa.Port
	http.ListenAndServe(addr, nil)
}

func (ioa *Ioa) ReverseProxy(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	path := r.URL.Path
	apiId, params, _ := ioa.Router.FindRoute(method, path)
	log.Println("api id is :", apiId)
	log.Println("api params is :", params)
	//todo match api, get model.Api
	api := ioa.Apis[apiId]

	for _, pluginId := range api.Plugins {
		name := ioa.Plugins[pluginId].GetName()
		log.Println("plugin will run", name)
		err := ioa.Plugins[pluginId].Run(w, r, &api)
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
	log.Println("apissssss", ioa.Apis)
	for id, api := range ioa.Apis {
		ioa.Router.AddRoute(api.Method, api.Path, id)
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
			json.Unmarshal([]byte(api.Policies), &newApiPolicies)
			json.Unmarshal([]byte(api.Plugins), &newApiPlugins)

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

			var newApiPluginRawConfig map[string]string
			//处理 api 的 pluginConfig
			if api.PluginConfig != "" {
				err := json.Unmarshal([]byte(api.PluginConfig), &newApiPluginRawConfig)
				if err != nil {
					log.Println("api pluginConfig Unmarshal error")
				}
			}

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
				PluginRawConfig: newApiPluginRawConfig,
			}

			ioa.Apis[api.Id] = newApi
		}
	}
}
