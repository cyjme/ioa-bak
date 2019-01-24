package store

import (
	"github.com/coreos/etcd/clientv3"
	"time"
)

var client *clientv3.Client

const defaultTimeout = time.Second * 5

func Init() {
	var err error
	client, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Error(ERR_STORE_CREATE_ETCD_CLIENT, err)
		panic(err)
	}
}
