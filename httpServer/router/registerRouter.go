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
	}

	pluginController := controller.PluginController{}
	r.GET("/plugins", pluginController.List)
	r.GET("/pluginsWithTag", pluginController.ListWithTag)
}
