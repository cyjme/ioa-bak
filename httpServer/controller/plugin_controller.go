//generate by gen
package controller

import (
	"github.com/gin-gonic/gin"
	"ioa"
	"ioa/httpServer/model"
	"ioa/httpServer/pkg"
	"net/http"
	"strconv"
)

type PluginController struct {
}

// @Summary Create
// @Tags    Plugin
// @Param body body model.Plugin true "Plugin"
// @Success 200 {string} json ""
// @Router /plugins [post]
func (ctl *PluginController) Create(c *gin.Context) {
	plugin := model.Plugin{}

	if err := pkg.ParseRequest(c, &plugin); err != nil {
		return
	}

	if err := plugin.Insert(); err != nil {
		c.JSON(http.StatusBadGateway, err)
		return
	}

	c.JSON(http.StatusOK, plugin)
}

// @Summary  Delete
// @Tags     Plugin
// @Param  pluginId  path string true "pluginId"
// @Success 200 {string} json ""
// @Router /plugins/{pluginId} [delete]
func (ctl *PluginController) Delete(c *gin.Context) {
	plugin := model.Plugin{}
	plugin.Id = c.Param("pluginId")
	err := plugin.Delete()

	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, nil)
}

// @Summary Put
// @Tags    Plugin
// @Param body body model.Plugin true "plugin"
// @Param  pluginId path string true "pluginId"
// @Success 200 {string} json ""
// @Router /plugins/{pluginId} [put]
func (ctl *PluginController) Put(c *gin.Context) {
	plugin := model.Plugin{}
	plugin.Id = c.Param("pluginId")

	if err := pkg.ParseRequest(c, &plugin); err != nil {
		return
	}

	err := plugin.Update()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, nil)
}

// @Summary Patch
// @Tags    Plugin
// @Param body body model.Plugin true "plugin"
// @Param  pluginId path string true "pluginId"
// @Success 200 {string} json ""
// @Router /plugins/{pluginId} [patch]
func (ctl *PluginController) Patch(c *gin.Context) {
	plugin := model.Plugin{}
	plugin.Id = c.Param("pluginId")

	if err := pkg.ParseRequest(c, &plugin); err != nil {
		return
	}

	err := plugin.Patch()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, nil)
}

// @Summary List
// @Tags    Plugin
// @Param query query string false "query, ?query=age:>:21,name:like:%jason%"
// @Param order query string false "order, ?order=age:desc,created_at:asc"
// @Param page query int false "page"
// @Param pageSize query int false "pageSize"
// @Success 200 {array} model.Plugin "plugin array"
// @Router /plugins [get]
func (ctl *PluginController) List(c *gin.Context) {
	plugin := &model.Plugin{}
	plugin.Id = c.Param("pluginId")
	var err error

	pageParam := c.DefaultQuery("page", "-1")
	pageSizeParam := c.DefaultQuery("pageSize", "-1")
	rawQuery := c.DefaultQuery("query", "")
	rawOrder := c.DefaultQuery("order", "")

	pageInt, err := strconv.Atoi(pageParam)
	pageSizeInt, err := strconv.Atoi(pageSizeParam)

	offset := pageInt*pageSizeInt - pageSizeInt
	limit := pageSizeInt

	if pageInt < 0 || pageSizeInt < 0 {
		limit = -1
	}

	plugins, total, err := plugin.List(rawQuery, rawOrder, offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total": total,
		"data":  plugins,
	})
}

// @Summary Get
// @Tags    Plugin
// @Param  pluginId path string true "pluginId"
// @Success 200 {object} model.Plugin "plugin object"
// @Router /plugins/{pluginId} [get]
func (ctl *PluginController) Get(c *gin.Context) {
	plugin := &model.Plugin{}
	plugin.Id = c.Param("pluginId")

	plugin, err := plugin.Get()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, plugin)
}

// @Summary Get
// @Tags    Plugin
// @Param  pluginId path string true "pluginId"
// @Success 200 {array} plugin.Field "plugin ConfigTpl"
// @Router /plugins/{pluginId}/configTpl [get]
func (ctl *PluginController) GetPluginConfigTpl(c *gin.Context,ioa *ioa.Ioa) {
	configTpl := ioa.Plugins.GetPluginConfigTpl(c.Param("pluginId"))
	c.JSON(http.StatusOK, configTpl)
}
