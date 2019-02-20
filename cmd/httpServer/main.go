package main

import (
	"flag"
	"github.com/coreos/etcd/clientv3"
	"github.com/gin-gonic/gin"
	"ioa"
	"ioa/httpServer/pkg/middleware"
	"ioa/httpServer/router"
	logger "ioa/log"
	"ioa/store"
	_ "net/http/pprof"
	"time"
)

var log = logger.Get()

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

	r := gin.Default()
	r.Use(middleware.CORSMiddleware())
	router.RegisterRouter(r)

	addr := config.HttpServer.Host + ":" + config.HttpServer.Port

	log.Info("httpServer run at: ", addr)

	err := r.Run(addr)
	if err != nil {
		panic(err)
	}
}
