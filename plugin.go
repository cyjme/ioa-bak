package ioa

import (
	"github.com/sirupsen/logrus"
	"ioa/proto"
	"net/http"
	goPlugin "plugin"
)

type Plugins map[string]IoaPlugin

type IoaPlugin interface {
	GetName() string
	GetDescribe() string
	GetConfigTemplate() proto.ConfigTpl
	InitApi(api *Api) error
	Run(w http.ResponseWriter, r *http.Request, api *Api) error
}

type Plugin struct {
	Name      string          `json:"name"`
	Describe  string          `json:"describe"`
	ConfigTpl proto.ConfigTpl `json:"configTpl"`
}

func (p Plugins) GetPluginConfigTpl(id string) proto.ConfigTpl {
	var configTpl proto.ConfigTpl
	configTpl = p[id].GetConfigTemplate()

	return configTpl
}

func (p Plugins) Register(id string, path string) {
	plugin, err := goPlugin.Open(path)

	if err != nil {
		log.Debug(err.Error())
	}

	symbol, err := plugin.Lookup("IoaPlugin")
	if err != nil {
		log.Debug("lookup plugin error", err.Error())
	}

	var ioaPlugin IoaPlugin
	ioaPlugin, ok := symbol.(IoaPlugin)

	if !ok {
		log.Debug("load plugin error")
		return
	}

	p[id] = ioaPlugin
}

type BasePlugin struct {
}

func (b *BasePlugin) Logger() *logrus.Logger {
	return log
}
