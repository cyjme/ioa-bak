package plugin

import (
	"log"
	"net/http"
	"plugin"
)

type PluginCenter struct {
	Plugins map[string]IoaPlugin
}

type IoaPlugin interface {
	GetName() string
	GetConfig()
	Run(w http.ResponseWriter, r *http.Request)
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

	symbol, err := plugin.Lookup("SizePlugin")
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
