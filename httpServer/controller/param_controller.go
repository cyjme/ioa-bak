//generate by gen
package controller

import (
	"github.com/gin-gonic/gin"
	"ioa/httpServer/model"
	"ioa/httpServer/pkg"
	"net/http"
	"strconv"
)

type ParamController struct {
}

// @Summary Create
// @Tags    Param
// @Param body body model.Param true "Param"
// @Success 200 {string} json ""
// @Router /params [post]
func (ctl *ParamController) Create(c *gin.Context) {
	param := model.Param{}

	if err := pkg.ParseRequest(c, &param); err != nil {
		return
	}

	if err := param.Insert(); err != nil {
		c.JSON(http.StatusBadGateway, err)
		return
	}

	c.JSON(http.StatusOK, param)
}

// @Summary  Delete
// @Tags     Param
// @Param  paramId  path string true "paramId"
// @Success 200 {string} json ""
// @Router /params/{paramId} [delete]
func (ctl *ParamController) Delete(c *gin.Context) {
	param := model.Param{}
	param.Id = c.Param("paramId")
	err := param.Delete()

	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, nil)
}

// @Summary Put
// @Tags    Param
// @Param body body model.Param true "param"
// @Param  paramId path string true "paramId"
// @Success 200 {string} json ""
// @Router /params/{paramId} [put]
func (ctl *ParamController) Put(c *gin.Context) {
	param := model.Param{}
	param.Id = c.Param("paramId")

	if err := pkg.ParseRequest(c, &param); err != nil {
		return
	}

	err := param.Update()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, nil)
}

// @Summary Patch
// @Tags    Param
// @Param body body model.Param true "param"
// @Param  paramId path string true "paramId"
// @Success 200 {string} json ""
// @Router /params/{paramId} [patch]
func (ctl *ParamController) Patch(c *gin.Context) {
	param := model.Param{}
	param.Id = c.Param("paramId")

	if err := pkg.ParseRequest(c, &param); err != nil {
		return
	}

	err := param.Patch()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, nil)
}

// @Summary List
// @Tags    Param
// @Param query query string false "query, ?query=age:>:21,name:like:%jason%"
// @Param order query string false "order, ?order=age:desc,created_at:asc"
// @Param page query int false "page"
// @Param pageSize query int false "pageSize"
// @Success 200 {array} model.Param "param array"
// @Router /params [get]
func (ctl *ParamController) List(c *gin.Context) {
	param := &model.Param{}
	param.Id = c.Param("paramId")
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

	params, total, err := param.List(rawQuery, rawOrder, offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total": total,
		"data":  params,
	})
}

// @Summary Get
// @Tags    Param
// @Param  paramId path string true "paramId"
// @Success 200 {object} model.Param "param object"
// @Router /params/{paramId} [get]
func (ctl *ParamController) Get(c *gin.Context) {
	param := &model.Param{}
	param.Id = c.Param("paramId")

	param, err := param.Get()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, param)
}
