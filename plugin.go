package ioa

import (
	"github.com/sirupsen/logrus"
	"ioa/proto"
	goPlugin "plugin"
)

type Plugins map[string]Plugin

type Plugin interface {
	GetName() string
	GetDescribe() string
	GetConfigTemplate() proto.ConfigTpl
	InitApi(api *Api) error
	ReceiveRequest(context *Context)
	ReceiveResponse(context *Context)
}

func (p Plugins) GetPluginConfigTpl(id string) proto.ConfigTpl {
	var configTpl proto.ConfigTpl
	configTpl = p[id].GetConfigTemplate()

	return configTpl
}

func (p Plugins) Register(id string, path string) {
	plugin, err := goPlugin.Open(path)

	if err != nil {
		log.Error(ERR_PLUGIN_OPEN_FILE, err.Error())
	}

	symbol, err := plugin.Lookup("ExportPlugin")
	if err != nil {
		log.Error(ERR_PLUGIN_LOOKUP, err.Error())
	}

	var ioaPlugin Plugin
	ioaPlugin, ok := symbol.(Plugin)

	if !ok {
		log.Error(ERR_PLUGIN_TYPE_ASSERTION)
		return
	}

	p[id] = ioaPlugin
}

type BasePlugin struct {
}

func (b *BasePlugin) Logger() *logrus.Logger {
	return log
}
