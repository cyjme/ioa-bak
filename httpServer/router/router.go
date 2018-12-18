package router

import (
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	_ "ioa/docs"
	"ioa/httpServer/app"
	"ioa/httpServer/controller"
	"ioa/httpServer/pkg/middleware"
)

func Start() {

	r := gin.Default()
	r.Use(middleware.CORSMiddleware())
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
		apiGroup.PATCH("/:apiId", apiController.Patch)
	}

	apiGroupController := controller.ApiGroupController{}
	apiGroupGroup := r.Group("/apiGroups")
	{
		apiGroupGroup.GET("", apiGroupController.List)
		apiGroupGroup.POST("", apiGroupController.Create)
		apiGroupGroup.DELETE("/:apiGroupId", apiGroupController.Delete)
		apiGroupGroup.PUT("/:apiGroupId", apiGroupController.Put)
		apiGroupGroup.GET("/:apiGroupId", apiGroupController.Get)
		apiGroupGroup.PATCH("/:apiGroupId", apiGroupController.Patch)
	}

	pluginController := controller.PluginController{}
	pluginGroup := r.Group("/plugins")
	{
		pluginGroup.GET("", pluginController.List)
		pluginGroup.POST("", pluginController.Create)
		pluginGroup.DELETE("/:pluginId", pluginController.Delete)
		pluginGroup.PUT("/:pluginId", pluginController.Put)
		pluginGroup.GET("/:pluginId", pluginController.Get)
		pluginGroup.PATCH("/:pluginId", pluginController.Patch)
	}

	policyController := controller.PolicyController{}
	policyGroup := r.Group("/policys")
	{
		policyGroup.GET("", policyController.List)
		policyGroup.POST("", policyController.Create)
		policyGroup.DELETE("/:policyId", policyController.Delete)
		policyGroup.PUT("/:policyId", policyController.Put)
		policyGroup.GET("/:policyId", policyController.Get)
		policyGroup.PATCH("/:policyId", policyController.Patch)
	}

	//!!do not delete gen will generate router code at here

	addr := app.Config.Http.Domain + ":" + app.Config.Http.Port

	r.Run(addr) // listen and serve on 0.0.0.0:8080
}
