//generate by gen
package controller

import (
	"github.com/gin-gonic/gin"
	"httpServer/model"
	"httpServer/pkg"
	"net/http"
	"strconv"
)

type ApiGroupController struct {
}

// @Summary Create
// @Tags    ApiGroup
// @Param body body model.ApiGroup true "ApiGroup"
// @Success 200 {string} json ""
// @Router /apiGroups [post]
func (ctl *ApiGroupController) Create(c *gin.Context) {
	apiGroup := model.ApiGroup{}

	if err := pkg.ParseRequest(c, &apiGroup); err != nil {
		return
	}

	if err := apiGroup.Insert(); err != nil {
		c.JSON(http.StatusBadGateway, err)
		return
	}

	c.JSON(http.StatusOK, apiGroup)
}

// @Summary  Delete
// @Tags     ApiGroup
// @Param  apiGroupId  path string true "apiGroupId"
// @Success 200 {string} json ""
// @Router /apiGroups/{apiGroupId} [delete]
func (ctl *ApiGroupController) Delete(c *gin.Context) {
	apiGroup := model.ApiGroup{}
	apiGroup.Id = c.Param("apiGroupId")
	err := apiGroup.Delete()

	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, nil)
}

// @Summary Put
// @Tags    ApiGroup
// @Param body body model.ApiGroup true "apiGroup"
// @Param  apiGroupId path string true "apiGroupId"
// @Success 200 {string} json ""
// @Router /apiGroups/{apiGroupId} [put]
func (ctl *ApiGroupController) Put(c *gin.Context) {
	apiGroup := model.ApiGroup{}
	apiGroup.Id = c.Param("apiGroupId")

	if err := pkg.ParseRequest(c, &apiGroup); err != nil {
		return
	}

	err := apiGroup.Update()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, nil)
}

// @Summary Patch
// @Tags    ApiGroup
// @Param body body model.ApiGroup true "apiGroup"
// @Param  apiGroupId path string true "apiGroupId"
// @Success 200 {string} json ""
// @Router /apiGroups/{apiGroupId} [patch]
func (ctl *ApiGroupController) Patch(c *gin.Context) {
	apiGroup := model.ApiGroup{}
	apiGroup.Id = c.Param("apiGroupId")

	if err := pkg.ParseRequest(c, &apiGroup); err != nil {
		return
	}

	err := apiGroup.Patch()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, nil)
}

// @Summary List
// @Tags    ApiGroup
// @Param query query string false "query, ?query=age:>:21,name:like:%jason%"
// @Param order query string false "order, ?order=age:desc,created_at:asc"
// @Param page query int false "page"
// @Param pageSize query int false "pageSize"
// @Success 200 {array} model.ApiGroup "apiGroup array"
// @Router /apiGroups [get]
func (ctl *ApiGroupController) List(c *gin.Context) {
	apiGroup := &model.ApiGroup{}
	apiGroup.Id = c.Param("apiGroupId")
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

	apiGroups, total, err := apiGroup.List(rawQuery, rawOrder, offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total": total,
		"data":  apiGroups,
	})
}

// @Summary Get
// @Tags    ApiGroup
// @Param  apiGroupId path string true "apiGroupId"
// @Success 200 {object} model.ApiGroup "apiGroup object"
// @Router /apiGroups/{apiGroupId} [get]
func (ctl *ApiGroupController) Get(c *gin.Context) {
	apiGroup := &model.ApiGroup{}
	apiGroup.Id = c.Param("apiGroupId")

	apiGroup, err := apiGroup.Get()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, apiGroup)
}
