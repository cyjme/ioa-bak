package ioa

import (
	"encoding/json"
	"ioa/store"
)

type Api struct {
	ApiGroupId string `json:"apiGroupId"`
	Name       string `json:"name"`
	Describe   string `json:"describe"`
	Path       string `json:"path"`
	Method     string `json:"method"`
	Status     string `json:"status"`

	Targets         []store.Target             `json:"targets"`
	Plugins         []string                   `json:"plugins"`

	PluginRawConfig map[string]json.RawMessage `json:"pluginRawConfig"`
	PluginConfig    map[string]interface{}     `json:"pluginConfig"`
	PluginsData     map[string]interface{}     `json:"pluginsData"`
}
