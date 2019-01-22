package ioa

import (
	"github.com/sirupsen/logrus"
	"ioa/proto"
	"net/http"
	goPlugin "plugin"
)

type Plugins map[string]Plugin

type Plugin interface {
	GetName() string
	GetDescribe() string
	GetConfigTemplate() proto.ConfigTpl
	InitApi(api *Api) error
	ReceiveRequest(context *Context) error
	ReceiveResponse(context *Context) error
}

type Context struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request
	Response       *http.Response
	Api            *Api
	Cancel         bool
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

	symbol, err := plugin.Lookup("ExportPlugin")
	if err != nil {
		log.Debug("lookup plugin error", err.Error())
	}

	var ioaPlugin Plugin
	ioaPlugin, ok := symbol.(Plugin)

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
