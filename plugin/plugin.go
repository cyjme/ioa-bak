package plugin

import (
	"log"
	"net/http"
	"plugin"
)

type PluginCenter struct {
	Plugins map[string]IoaPlugin
}

type Field struct {
	Name      string
	Required  bool
	FieldType string
}

type Config []Field

type IoaPlugin interface {
	GetName() string
	GetConfigTemplate() Config
	Run(w http.ResponseWriter, r *http.Request, config map[string]interface{})
}

func NewPluginCenter() *PluginCenter {
	return &PluginCenter{
		Plugins: make(map[string]IoaPlugin),
	}
}

func (p *PluginCenter) Register(id string, path string) {
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

	p.Plugins[id] = ioaPlugin
}
