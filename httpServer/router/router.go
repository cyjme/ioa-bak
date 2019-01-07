package router

import (
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"ioa"
	_ "ioa/docs"
	"ioa/httpServer/app"
	"ioa/httpServer/controller"
	"ioa/httpServer/pkg/middleware"
)

func Start(ioa *ioa.Ioa) {

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

	r.GET("/test/users", func(context *gin.Context) {
		type user struct {
			Name  string `json:"name"`
			Phone string `json:"phone"`
			Email string `json:"email"`
		}
		data := []user{
			user{
				Name:  "jason",
				Phone: "13061710381",
				Email: "changyuanjian@gmail.com",
			},
			user{
				Name:  "wwq",
				Phone: "13023221023",
				Email: "wwq@kuipgroup.com",
			},
		}
		context.JSON(200, data)
	})
	r.POST("/test/users", func(context *gin.Context) {
		context.JSON(200, nil)
	})

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
		pluginGroup.GET("", func(c *gin.Context) {
			pluginController.List(c, ioa)
		})

		pluginGroup.GET("/:pluginName", func(c *gin.Context) {
			pluginController.Get(c, ioa)
		})
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

	paramController := controller.ParamController{}
	paramGroup := r.Group("/params")
	{
		paramGroup.GET("", paramController.List)
		paramGroup.POST("", paramController.Create)
		paramGroup.DELETE("/:paramId", paramController.Delete)
		paramGroup.PUT("/:paramId", paramController.Put)
		paramGroup.GET("/:paramId", paramController.Get)
		paramGroup.PATCH("/:paramId", paramController.Patch)
	}

	targetController := controller.TargetController{}
	targetGroup := r.Group("/targets")
	{
		targetGroup.GET("", targetController.List)
		targetGroup.POST("", targetController.Create)
		targetGroup.DELETE("/:targetId", targetController.Delete)
		targetGroup.PUT("/:targetId", targetController.Put)
		targetGroup.GET("/:targetId", targetController.Get)
		targetGroup.PATCH("/:targetId", targetController.Patch)
	}

	//!!do not delete gen will generate router code at here

	addr := app.Config.Http.Host + ":" + app.Config.Http.Port

	r.Run(addr) // listen and serve on 0.0.0.0:8080
}
