package controller

import (
	"github.com/gin-gonic/gin"
	"ioa/proto"
	"ioa/store"
	"net/http"
)

type PluginController struct {
}

// @Summary Get
// @Tags   Plugin
// @Router /plugins [get]
func (p *PluginController) List(c *gin.Context) {
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

// @Summary Get
// @Tags   Plugin
// @Router /pluginsWithTag [get]
func (p *PluginController) ListWithTag(c *gin.Context) {
	plugins, total, err := store.ListPlugin()
	if err != nil {
		c.JSON(http.StatusBadGateway, err)
		return
	}

	tag2plugins := make(map[string][]proto.Plugin, 0)

	for _, plugin := range plugins {
		for _, tag := range plugin.Tags {
			oldTagPlugins := tag2plugins[tag]
			tag2plugins[tag] = append(oldTagPlugins, plugin)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"total": total,
		"data":  tag2plugins,
	})
}
