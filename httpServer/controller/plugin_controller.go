package controller

import (
	"github.com/gin-gonic/gin"
	"ioa/store"
	"net/http"
)

type PluginController struct {
}

// @Summary Get
// @Tags   Plugin
// @Router /plugins [get]
func (p *PluginController) List(c *gin.Context)  {
	plugins, total, err := store.ListPlugin()
	if err != nil {
		c.JSON(http.StatusBadGateway, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total": total,
		"data":  plugins,
	})
}
