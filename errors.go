package ioa

import "errors"

var (
	ERR_INIT_API_PLUGIN        = errors.New("run_plugin: init api plugin error\n")
	ERR_PLUGIN_TYPE_ASSERTION  = errors.New("plugin: type assertion error\n")
	ERR_PLUGIN_OPEN_FILE       = errors.New("plugin: open .so file error\n")
	ERR_PLUGIN_LOOKUP          = errors.New("plugin: lookup Plugin error")
	ERR_PROXY_CREATE_REQUEST   = errors.New("proxy: create request error")
	ERR_PROXY_DO_REQUEST       = errors.New("proxy: do request error")
	ERR_API_GET_PLUGINS        = errors.New("api: get plugin by Unmarshal api.plugins")
	ERR_API_USE_UNEXIST_PLUGIN = errors.New("api: use unexist plugin")
)
