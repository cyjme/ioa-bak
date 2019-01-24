package store

import "errors"

var (
	ERR_STORE_CRUD_API           = errors.New("ERR_STORE_CRUD_API")
	ERR_STORE_CREATE_PLUGIN      = errors.New("ERR_STORE_CREATE_PLUGIN")
	ERR_STORE_LIST_PLUGIN        = errors.New("ERR_STORE_LIST_PLUGIN")
	ERR_STORE_CREATE_ETCD_CLIENT = errors.New("ERR_STORE_CREATE_ETCD_CLIENT")
)
