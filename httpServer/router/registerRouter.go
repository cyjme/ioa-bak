package router

import (
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"ioa/httpServer/controller"
)

func RegisterRouter(r *gin.Engine) {
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	if gin.Mode() == gin.DebugMode {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	apiController := controller.ApiController{}
	apiGroup := r.Group("/apis")
	{
		apiGroup.GET("", apiController.List)
		apiGroup.POST("", apiController.Create)
		apiGroup.DELETE("/:apiId", apiController.Delete)
		apiGroup.PUT("/:apiId", apiController.Put)
		apiGroup.GET("/:apiId", apiController.Get)
	}

	r.GET("/apisWithTag", apiController.ListWithTag)

	pluginController := controller.PluginController{}
	r.GET("/plugins", pluginController.List)
	r.GET("/pluginsWithTag", pluginController.ListWithTag)
}
