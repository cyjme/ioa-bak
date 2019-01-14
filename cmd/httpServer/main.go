package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"ioa/httpServer/pkg/middleware"
	"ioa/httpServer/router"
	"ioa/store"
	"log"
	_ "net/http/pprof"
)

func main() {
	store.Init()
	r := gin.Default()
	r.Use(middleware.CORSMiddleware())
	router.RegisterRouter(r)

	var addr string
	flag.StringVar(&addr, "addr", "0.0.0.0:9992", "")
	flag.Parse()
	log.Println("httpServer run at: ", addr)

	err := r.Run(addr)
	if err != nil {
		panic(err)
	}
}
