package router

import (
	"github.com/gin-gonic/gin"
	_ "ioa/docs"
	"ioa/httpServer/app"
	"ioa/httpServer/pkg/middleware"
)

func Start() {
	r := gin.Default()
	r.Use(middleware.CORSMiddleware())

	//!!do not delete gen will generate router code at here

	addr := app.Config.Http.Host + ":" + app.Config.Http.Port

	err := r.Run(addr) // listen and serve on 0.0.0.0:8080

	if err != nil {
		panic(err)
	}
}
