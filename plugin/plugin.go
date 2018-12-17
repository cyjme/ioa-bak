package plugin

import (
	"log"
	"net/http"
	"plugin"
)

type PluginCenter struct {
	Plugins map[string]Plugin
}

type Plugin struct {
	Id  string
	Run func(w http.ResponseWriter, r *http.Request)
}

func NewPluginCenter() *PluginCenter {
	return &PluginCenter{
		Plugins: make(map[string]Plugin),
	}
}

func (p *PluginCenter) Register(id string, path string) {
	plugin, err := plugin.Open(path)
	if err != nil {
		log.Println(err.Error())
	}

	symbol, err := plugin.Lookup("Run")
	if err != nil {
		log.Println(err.Error())
	}
	runFunction, ok := symbol.(func(w http.ResponseWriter, r *http.Request))

	if !ok {
		log.Println("function type not match func(w http.ResponseWriter, r *http.Request)")
		return
	}

	p.Plugins[id] = Plugin{
		Id:  id,
		Run: runFunction,
	}
}
