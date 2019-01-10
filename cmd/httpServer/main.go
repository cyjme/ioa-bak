package main

import (
	"github.com/gin-gonic/gin"
	"ioa/httpServer/pkg/middleware"
	"ioa/httpServer/router"
	"ioa/store"
	_ "net/http/pprof"
)

func main() {
	store.Init()
	r := gin.Default()
	r.Use(middleware.CORSMiddleware())
	router.RegisterRouter(r)
	addr := "0.0.0.0:9992"

	err := r.Run(addr) // listen and serve on 0.0.0.0:8080
	if err != nil {
		panic(err)
	}
}
