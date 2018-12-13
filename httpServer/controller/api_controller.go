//generate by gen
package controller

import (
	"github.com/gin-gonic/gin"
	"httpServer/model"
	"httpServer/pkg"
	"net/http"
	"strconv"
)

type ApiController struct {
}

// @Summary Create
// @Tags    Api
// @Param body body model.Api true "Api"
// @Success 200 {string} json ""
// @Router /apis [post]
func (ctl *ApiController) Create(c *gin.Context) {
	api := model.Api{}

	if err := pkg.ParseRequest(c, &api); err != nil {
		return
	}

	if err := api.Insert(); err != nil {
		c.JSON(http.StatusBadGateway, err)
		return
	}

	c.JSON(http.StatusOK, api)
}

// @Summary  Delete
// @Tags     Api
// @Param  apiId  path string true "apiId"
// @Success 200 {string} json ""
// @Router /apis/{apiId} [delete]
func (ctl *ApiController) Delete(c *gin.Context) {
	api := model.Api{}
	api.Id = c.Param("apiId")
	err := api.Delete()

	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, nil)
}

// @Summary Put
// @Tags    Api
// @Param body body model.Api true "api"
// @Param  apiId path string true "apiId"
// @Success 200 {string} json ""
// @Router /apis/{apiId} [put]
func (ctl *ApiController) Put(c *gin.Context) {
	api := model.Api{}
	api.Id = c.Param("apiId")

	if err := pkg.ParseRequest(c, &api); err != nil {
		return
	}

	err := api.Update()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, nil)
}

// @Summary Patch
// @Tags    Api
// @Param body body model.Api true "api"
// @Param  apiId path string true "apiId"
// @Success 200 {string} json ""
// @Router /apis/{apiId} [patch]
func (ctl *ApiController) Patch(c *gin.Context) {
	api := model.Api{}
	api.Id = c.Param("apiId")

	if err := pkg.ParseRequest(c, &api); err != nil {
		return
	}

	err := api.Patch()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, nil)
}

// @Summary List
// @Tags    Api
// @Param query query string false "query, ?query=age:>:21,name:like:%jason%"
// @Param order query string false "order, ?order=age:desc,created_at:asc"
// @Param page query int false "page"
// @Param pageSize query int false "pageSize"
// @Success 200 {array} model.Api "api array"
// @Router /apis [get]
func (ctl *ApiController) List(c *gin.Context) {
	api := &model.Api{}
	api.Id = c.Param("apiId")
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

	apis, total, err := api.List(rawQuery, rawOrder, offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total": total,
		"data":  apis,
	})
}

// @Summary Get
// @Tags    Api
// @Param  apiId path string true "apiId"
// @Success 200 {object} model.Api "api object"
// @Router /apis/{apiId} [get]
func (ctl *ApiController) Get(c *gin.Context) {
	api := &model.Api{}
	api.Id = c.Param("apiId")

	api, err := api.Get()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, api)
}
