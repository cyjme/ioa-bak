package ioa

import (
	"encoding/json"
	"ioa/httpServer/app"
	"ioa/httpServer/model"
	"ioa/plugin"
	"ioa/router"
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
)

type Ioa struct {
	Apis    map[string]Api
	Plugins plugin.Plugin
	Router  router.Router
}

func New() *Ioa {
	return &Ioa{
		Apis:    make(map[string]Api),
		Plugins: make(plugin.Plugin),
		Router:  router.New(),
	}
}

type Api struct {
	model.Common
	ApiGroupId string `json:"apiGroupId"`
	Name       string `json:"name"`
	Describe   string `json:"describe"`
	Path       string `json:"path"`
	Method     string `json:"method"`
	Status     string `json:"status"`

	Targets       []model.Target    `json:"targets"`
	Params        []model.Param     `json:"params"`
	Policies      []string          `json:"policies"`
	GroupPolicies []string          `json:"groupPolicies"`
	Plugins       []string          `json:"plugins"`
	GroupPlugins  []string          `json:"groupPlugins"`
	AllPlugin     []string          `json:"allPlugin"`
	PluginConfig  map[string]string `json:"pluginConfig"`
}

func (ioa *Ioa) StartServer() {
	log.Println("load Api from database")
	ioa.LoadApi()
	ioa.loadApiToRouter()

	log.Println("load Api from database Success:", ioa.Apis)
	ioa.Plugins.Register("1", "./plugins/size.so")

	http.HandleFunc("/", ioa.ReverseProxy)

	addr := app.Config.Ioa.Host + ":" + app.Config.Ioa.Port
	http.ListenAndServe(addr, nil)
}

func (ioa *Ioa) ReverseProxy(w http.ResponseWriter, r *http.Request) {
	log.Println("***********************", r.URL.Scheme)
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
		ioa.Plugins[pluginId].Run(w, r, map[string]interface{}{"maxSize": int64(10000)})
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

			var newApiPluginConfig map[string]string
			//处理 api 的 pluginConfig
			if api.PluginConfig != "" {
				err := json.Unmarshal([]byte(api.PluginConfig), &newApiPluginConfig)
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

				Targets:       api.Targets,
				Params:        api.Params,
				Policies:      newApiPolicies,
				GroupPolicies: apiGroupPolicies,
				GroupPlugins:  apiGroupPlugins,
				Plugins:       newApiPlugins,
				AllPlugin:     newApiAllPlugins,
				PluginConfig:  newApiPluginConfig,
			}

			ioa.Apis[api.Id] = newApi
		}
	}
}
