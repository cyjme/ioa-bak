//generate by gen
package controller

import (
	"github.com/gin-gonic/gin"
	"ioa/httpServer/model"
	"ioa/httpServer/pkg"
	"net/http"
	"strconv"
)

type TargetController struct {
}

// @Summary Create
// @Tags    Target
// @Param body body model.Target true "Target"
// @Success 200 {string} json ""
// @Router /targets [post]
func (ctl *TargetController) Create(c *gin.Context) {
	target := model.Target{}

	if err := pkg.ParseRequest(c, &target); err != nil {
		return
	}

	if err := target.Insert(); err != nil {
		c.JSON(http.StatusBadGateway, err)
		return
	}

	c.JSON(http.StatusOK, target)
}

// @Summary  Delete
// @Tags     Target
// @Param  targetId  path string true "targetId"
// @Success 200 {string} json ""
// @Router /targets/{targetId} [delete]
func (ctl *TargetController) Delete(c *gin.Context) {
	target := model.Target{}
	target.Id = c.Param("targetId")
	err := target.Delete()

	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, nil)
}

// @Summary Put
// @Tags    Target
// @Param body body model.Target true "target"
// @Param  targetId path string true "targetId"
// @Success 200 {string} json ""
// @Router /targets/{targetId} [put]
func (ctl *TargetController) Put(c *gin.Context) {
	target := model.Target{}
	target.Id = c.Param("targetId")

	if err := pkg.ParseRequest(c, &target); err != nil {
		return
	}

	err := target.Update()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, nil)
}

// @Summary Patch
// @Tags    Target
// @Param body body model.Target true "target"
// @Param  targetId path string true "targetId"
// @Success 200 {string} json ""
// @Router /targets/{targetId} [patch]
func (ctl *TargetController) Patch(c *gin.Context) {
	target := model.Target{}
	target.Id = c.Param("targetId")

	if err := pkg.ParseRequest(c, &target); err != nil {
		return
	}

	err := target.Patch()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, nil)
}

// @Summary List
// @Tags    Target
// @Param query query string false "query, ?query=age:>:21,name:like:%jason%"
// @Param order query string false "order, ?order=age:desc,created_at:asc"
// @Param page query int false "page"
// @Param pageSize query int false "pageSize"
// @Success 200 {array} model.Target "target array"
// @Router /targets [get]
func (ctl *TargetController) List(c *gin.Context) {
	target := &model.Target{}
	target.Id = c.Param("targetId")
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

	targets, total, err := target.List(rawQuery, rawOrder, offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total": total,
		"data":  targets,
	})
}

// @Summary Get
// @Tags    Target
// @Param  targetId path string true "targetId"
// @Success 200 {object} model.Target "target object"
// @Router /targets/{targetId} [get]
func (ctl *TargetController) Get(c *gin.Context) {
	target := &model.Target{}
	target.Id = c.Param("targetId")

	target, err := target.Get()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, target)
}
