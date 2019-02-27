package ioa

import "errors"

var (
	ERR_INIT_API_PLUGIN        = errors.New("ERR_INIT_API_PLUGIN")
	ERR_PLUGIN_TYPE_ASSERTION  = errors.New("ERR_PLUGIN_TYPE_ASSERTION")
	ERR_PLUGIN_OPEN_FILE       = errors.New("ERR_PLUGIN_OPEN_FILE")
	ERR_PLUGIN_LOOKUP          = errors.New("ERR_PLUGIN_LOOKUP")
	ERR_PROXY_CREATE_REQUEST   = errors.New("ERR_PROXY_CREATE_REQUEST")
	ERR_PROXY_DO_REQUEST       = errors.New("ERR_PROXY_DO_REQUEST")
	ERR_API_GET_PLUGINS        = errors.New("ERR_API_GET_PLUGINS")
	ERR_API_USE_UNEXIST_PLUGIN = errors.New("ERR_API_USE_UNEXIST_PLUGIN")
	ERR_CONFIG_LOAD            = errors.New("ERR_CONFIG_LOAD")
	ERR_ROUTER_ADD_EXISTED     = errors.New("ERR_ROUTER_ADD_EXISTED")
)
