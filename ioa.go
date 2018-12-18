package ioa

import (
	"encoding/json"
	"ioa/httpServer/model"
	"ioa/plugin"
	"log"
	"net/http"
)

var Count int64
var Apis []Api
var PluginCenter *plugin.PluginCenter

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
}

func StartServer() {
	log.Println("load Api from database")
	LoadApi()
	log.Println("load Api from database Success:", Apis)
	PluginCenter = plugin.NewPluginCenter()
	PluginCenter.Register("001", "./plugins/count.so")
	PluginCenter.Register("002", "./plugins/size.so")

	http.HandleFunc("/", ReverseProxy)
	http.ListenAndServe(":11112", nil)
}

func ReverseProxy(w http.ResponseWriter, r *http.Request) {
	PluginCenter.Plugins["001"].Run(w, r)
	PluginCenter.Plugins["002"].Run(w, r)
	w.Write([]byte("ok"))
	log.Println("receive request")
	//todo  use httpRouter get apiId

	//todo match api, get model.Api
	//todo plugin
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

			log.Println("test", plugins)

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
			}

			Apis = append(Apis, newApi)
		}
	}
}
