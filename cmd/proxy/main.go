package main

import (
	"flag"
	"github.com/coreos/etcd/clientv3"
	"ioa"
	"ioa/store"
	"time"
)

func main() {
	var path string
	flag.StringVar(&path, "config", "/etc/ioa", "")
	flag.Parse()

	config := ioa.ReadConfig(path)

	etcdConfig := clientv3.Config{
		Endpoints:   config.Etcd.Endpoints,
		DialTimeout: time.Duration(config.Etcd.DialTimeout) * time.Second,
		Username:    config.Etcd.Username,
		Password:    config.Etcd.Password,
	}

	store.Init(etcdConfig)
	Ioa := ioa.New(config)
	Ioa.StartServer()
}
