package ioa

import (
	"encoding/json"
	"ioa/httpServer/model"
)

type plugin struct {
	Name   string          `json:"name"`
	Config json.RawMessage `json:"config"`
}
type Api struct {
	model.Common
	ApiGroupId string `json:"apiGroupId"`
	Name       string `json:"name"`
	Describe   string `json:"describe"`
	Path       string `json:"path"`
	Method     string `json:"method"`
	Status     string `json:"status"`

	Targets         []model.Target             `json:"targets"`
	Params          []model.Param              `json:"params"`
	Policies        []string                   `json:"policies"`
	GroupPolicies   []string                   `json:"groupPolicies"`
	Plugins         []string                   `json:"plugins"`
	GroupPlugins    []string                   `json:"groupPlugins"`
	AllPlugin       []string                   `json:"allPlugin"`
	PluginRawConfig map[string]json.RawMessage `json:"pluginRawConfig"`
	PluginConfig    map[string]interface{}     `json:"pluginConfig"`
	PluginsData     map[string]interface{}     `json:"pluginsData"`
}
