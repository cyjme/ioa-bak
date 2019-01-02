package ioa

import (
	"log"
	"net/http"
	"plugin"
)

type Plugins map[string]IoaPlugin

type Field struct {
	Name      string `json:"name"`
	Desc      string `json:"desc"`
	Required  bool   `json:"required"`
	FieldType string `json:"fieldType"`
}

type ConfigTpl []Field

type IoaPlugin interface {
	GetName() string
	GetConfigTemplate() ConfigTpl
	InitApi(api *Api) error
	Run(w http.ResponseWriter, r *http.Request, api *Api) error
}

func (p Plugins) GetPluginConfigTpl(id string) ConfigTpl {
	var configTpl ConfigTpl
	configTpl = p[id].GetConfigTemplate()

	return configTpl
}

func (p Plugins) Register(id string, path string) {
	plugin, err := plugin.Open(path)

	if err != nil {
		log.Println(err.Error())
	}

	symbol, err := plugin.Lookup("IoaPlugin")
	if err != nil {
		log.Println("lookup plugin error", err.Error())
	}

	var ioaPlugin IoaPlugin
	ioaPlugin, ok := symbol.(IoaPlugin)

	if !ok {
		log.Println("load plugin error")
		return
	}

	p[id] = ioaPlugin
}