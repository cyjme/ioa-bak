package store

import (
	"github.com/coreos/etcd/clientv3"
	logger "ioa/log"
	"time"
)

var client *clientv3.Client

const apiPrefix = "/apis/"
const policyPrefix = "/policies/"
const defaultTimeout = time.Second * 5
var log = logger.Get()

func Init(config clientv3.Config) {
	var err error
	client, err = clientv3.New(config)

	if err != nil {
		log.Error(ERR_STORE_CREATE_ETCD_CLIENT, err)
		panic(err)
	}
}
