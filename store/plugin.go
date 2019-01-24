package store

import (
	"context"
	"encoding/json"
	"github.com/coreos/etcd/clientv3"
	"ioa/proto"
)

func ReCreatePlugin(plugins []proto.Plugin) error {
	pluginsByte, err := json.Marshal(plugins)
	if err != nil {
		log.Error(ERR_STORE_CREATE_PLUGIN, err)
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	_, err = clientv3.NewKV(client).Put(ctx, "/plugins", string(pluginsByte))
	cancel()
	if err != nil {
		log.Error(ERR_STORE_CREATE_PLUGIN, err)
		return err
	}

	return nil
}

func ListPlugin() ([]proto.Plugin, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	res, err := clientv3.NewKV(client).Get(ctx, "/plugins", clientv3.WithPrefix())
	cancel()
	if err != nil {
		log.Error(ERR_STORE_LIST_PLUGIN, err)
		return nil, 0, err
	}
	plugins := make([]proto.Plugin, 0)

	err = json.Unmarshal(res.Kvs[0].Value, &plugins)
	if err != nil {
		log.Error(ERR_STORE_LIST_PLUGIN, err)
		return nil, 0, err
	}

	return plugins, len(plugins), nil
}
