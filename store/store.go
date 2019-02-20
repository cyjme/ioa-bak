package store

import (
	"github.com/coreos/etcd/clientv3"
	"time"
)

var client *clientv3.Client

const defaultTimeout = time.Second * 5

func Init(config clientv3.Config) {
	var err error
	client, err = clientv3.New(config)

	if err != nil {
		log.Error(ERR_STORE_CREATE_ETCD_CLIENT, err)
		panic(err)
	}
}
