package ioa

import (
	"encoding/json"
	"ioa/httpServer/model"
	"ioa/plugin"
	"ioa/router"
	"log"
	"net/http"
)

var Apis = make(map[string]Api)
var PluginCenter *plugin.PluginCenter

type PluginConfig map[string]string

var Router = router.NewRouter()

type Api struct {
	model.Common
	ApiGroupId string `json:"apiGroupId"`
	Name       string `json:"name"`
	Describe   string `json:"describe"`
	Path       string `json:"path"`
	Method     string `json:"method"`
	Status     string `json:"status"`

	Targets       []model.Target `json:"targets"`
	Params        []model.Param  `json:"params"`
	Policies      []string       `json:"policies"`
	GroupPolicies []string       `json:"groupPolicies"`
	Plugins       []string       `json:"plugins"`
	GroupPlugins  []string       `json:"groupPlugins"`
	AllPlugin     []string       `json:"allPlugin"`
	PluginConfig  PluginConfig   `json:"pluginConfig"`
}

func StartServer() {
	log.Println("load Api from database")
	LoadApi()
	loadApiToRouter()

	log.Println("load Api from database Success:", Apis)
	PluginCenter = plugin.NewPluginCenter()
	PluginCenter.Register("001", "./plugins/size.so")

	http.HandleFunc("/", ReverseProxy)
	http.ListenAndServe(":11112", nil)
}

func ReverseProxy(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	path := r.URL.Path
	apiId, params, _ := Router.FindRoute(method, path)
	log.Println("api id is :", apiId)
	log.Println("api params is :", params)
	//todo match api, get model.Api

	PluginCenter.Plugins["001"].Run(w, r, map[string]interface{}{"maxSize": int64(10000)})
	name := PluginCenter.Plugins["001"].GetName()
	log.Println("name is ", name)

	w.Write([]byte("ok"))
	log.Println("receive request")

	//todo plugin
}

func loadApiToRouter() {
	log.Println("apissssss",Apis)
	for id, api := range Apis {
		Router.AddRoute(api.Method, api.Path, id)
	}
}

func LoadApi() {
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

			var newApiPluginConfig PluginConfig
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

			Apis[api.Id] = newApi
		}
	}
}
