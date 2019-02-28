//generate by gen
package controller

import (
	"github.com/gin-gonic/gin"
	"ioa/httpServer/pkg"
	"ioa/store"
	"ioa/util"
	"net/http"
)

type PolicyController struct {
}

// @Summary Create
// @Tags    Policy
// @Param body body store.Policy true "Policy"
// @Success 200 {string} json ""
// @Router /policies [post]
func (ctl *PolicyController) Create(c *gin.Context) {
	policy := store.Policy{}
	policy.Id = util.GetRandomString(11)

	if err := pkg.ParseRequest(c, &policy); err != nil {
		return
	}

	if err := policy.Put(); err != nil {
		c.JSON(http.StatusBadGateway, err)
		return
	}

	c.JSON(http.StatusOK, policy)
}

// @Summary  Delete
// @Tags     Policy
// @Param  policyId  path string true "policyId"
// @Success 200 {string} json ""
// @Router /policies/{policyId} [delete]
func (ctl *PolicyController) Delete(c *gin.Context) {
	policy := store.Policy{}
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
// @Param body body store.Policy true "policy"
// @Param  policyId path string true "policyId"
// @Success 200 {string} json ""
// @Router /policies/{policyId} [put]
func (ctl *PolicyController) Put(c *gin.Context) {
	policy := store.Policy{}
	policy.Id = c.Param("policyId")

	if err := pkg.ParseRequest(c, &policy); err != nil {
		return
	}

	err := policy.Put()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, nil)
}

// @Summary List
// @Tags    Policy
// @Success 200 {array} store.Policy "policy array"
// @Router /policies [get]
func (ctl *PolicyController) List(c *gin.Context) {
	policy := &store.Policy{}
	policy.Id = c.Param("policyId")
	var err error

	policies, total, err := policy.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total": total,
		"data":  policies,
	})
}

// @Summary Get
// @Tags    Policy
// @Param  policyId path string true "policyId"
// @Success 200 {object} store.Policy "policy object"
// @Router /policies/{policyId} [get]
func (ctl *PolicyController) Get(c *gin.Context) {
	policy := &store.Policy{}
	policy.Id = c.Param("policyId")

	policy, err := policy.Get()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, policy)
}
