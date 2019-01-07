package ioa

import (
	"log"
	"net/http"
	goPlugin "plugin"
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
	GetDescribe() string
	GetConfigTemplate() ConfigTpl
	InitApi(api *Api) error
	Run(w http.ResponseWriter, r *http.Request, api *Api) error
}

type Plugin struct {
	Name      string    `json:"name"`
	Describe  string    `json:"describe"`
	ConfigTpl ConfigTpl `json:"configTpl"`
}

func (p Plugins) GetPluginConfigTpl(id string) ConfigTpl {
	var configTpl ConfigTpl
	configTpl = p[id].GetConfigTemplate()

	return configTpl
}

func (p Plugins) Register(id string, path string) {
	plugin, err := goPlugin.Open(path)

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
