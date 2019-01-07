//generate by gen
package controller

import (
	"github.com/gin-gonic/gin"
	"ioa"
	"net/http"
)

type PluginController struct {
}

// @Summary List
// @Tags    Plugin
// @Router /plugins [get]
func (ctl *PluginController) List(c *gin.Context, i *ioa.Ioa) {
	var plugins []ioa.Plugin
	for _, plugin := range i.Plugins {
		plugins = append(plugins, ioa.Plugin{
			Name:      plugin.GetName(),
			Describe:  plugin.GetDescribe(),
			ConfigTpl: plugin.GetConfigTemplate(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"total": len(plugins),
		"data":  plugins,
	})
}

// @Summary GetPluginConfigTpl
// @Tags    Plugin
// @Param  pluginName path string true "pluginName"
// @Success 200 {array} ioa.Field "plugin ConfigTpl"
// @Router /plugins/{pluginName} [get]
func (ctl *PluginController) Get(c *gin.Context, i *ioa.Ioa) {
	configTpl := i.Plugins.GetPluginConfigTpl(c.Param("pluginName"))
	c.JSON(http.StatusOK, configTpl)
}
