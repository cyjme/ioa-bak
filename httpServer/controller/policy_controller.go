//generate by gen
package controller

import (
	"github.com/gin-gonic/gin"
	"ioa/httpServer/model"
	"ioa/httpServer/pkg"
	"net/http"
	"strconv"
)

type PolicyController struct {
}

// @Summary Create
// @Tags    Policy
// @Param body body model.Policy true "Policy"
// @Success 200 {string} json ""
// @Router /policys [post]
func (ctl *PolicyController) Create(c *gin.Context) {
	policy := model.Policy{}

	if err := pkg.ParseRequest(c, &policy); err != nil {
		return
	}

	if err := policy.Insert(); err != nil {
		c.JSON(http.StatusBadGateway, err)
		return
	}

	c.JSON(http.StatusOK, policy)
}

// @Summary  Delete
// @Tags     Policy
// @Param  policyId  path string true "policyId"
// @Success 200 {string} json ""
// @Router /policys/{policyId} [delete]
func (ctl *PolicyController) Delete(c *gin.Context) {
	policy := model.Policy{}
	policy.Id = c.Param("policyId")
	err := policy.Delete()

	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, nil)
}

// @Summary Put
// @Tags    Policy
// @Param body body model.Policy true "policy"
// @Param  policyId path string true "policyId"
// @Success 200 {string} json ""
// @Router /policys/{policyId} [put]
func (ctl *PolicyController) Put(c *gin.Context) {
	policy := model.Policy{}
	policy.Id = c.Param("policyId")

	if err := pkg.ParseRequest(c, &policy); err != nil {
		return
	}

	err := policy.Update()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, nil)
}

// @Summary Patch
// @Tags    Policy
// @Param body body model.Policy true "policy"
// @Param  policyId path string true "policyId"
// @Success 200 {string} json ""
// @Router /policys/{policyId} [patch]
func (ctl *PolicyController) Patch(c *gin.Context) {
	policy := model.Policy{}
	policy.Id = c.Param("policyId")

	if err := pkg.ParseRequest(c, &policy); err != nil {
		return
	}

	err := policy.Patch()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, nil)
}

// @Summary List
// @Tags    Policy
// @Param query query string false "query, ?query=age:>:21,name:like:%jason%"
// @Param order query string false "order, ?order=age:desc,created_at:asc"
// @Param page query int false "page"
// @Param pageSize query int false "pageSize"
// @Success 200 {array} model.Policy "policy array"
// @Router /policys [get]
func (ctl *PolicyController) List(c *gin.Context) {
	policy := &model.Policy{}
	policy.Id = c.Param("policyId")
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

	policys, total, err := policy.List(rawQuery, rawOrder, offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total": total,
		"data":  policys,
	})
}

// @Summary Get
// @Tags    Policy
// @Param  policyId path string true "policyId"
// @Success 200 {object} model.Policy "policy object"
// @Router /policys/{policyId} [get]
func (ctl *PolicyController) Get(c *gin.Context) {
	policy := &model.Policy{}
	policy.Id = c.Param("policyId")

	policy, err := policy.Get()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, policy)
}
