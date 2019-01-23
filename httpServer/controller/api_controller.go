//generate by gen
package controller

import (
	"github.com/gin-gonic/gin"
	"ioa/httpServer/pkg"
	"ioa/httpServer/pkg/util"
	"ioa/store"
	"net/http"
)

type ApiController struct {
}

// @Summary Create
// @Tags    Api
// @Param body body store.Api true "Api"
// @Success 200 {string} json ""
// @Router /apis [post]
func (ctl *ApiController) Create(c *gin.Context) {
	api := store.Api{}
	api.Id = util.GetRandomString(11)

	if err := pkg.ParseRequest(c, &api); err != nil {
		return
	}

	if err := api.Put(); err != nil {
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
	api := store.Api{}
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
// @Param body body store.Api true "api"
// @Param  apiId path string true "apiId"
// @Success 200 {string} json ""
// @Router /apis/{apiId} [put]
func (ctl *ApiController) Put(c *gin.Context) {
	api := store.Api{}
	api.Id = c.Param("apiId")

	if err := pkg.ParseRequest(c, &api); err != nil {
		return
	}

	err := api.Put()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, nil)
}

// @Summary ListWithTag
// @Tags    Api
// @Router /apisWithTag [get]
func (ctl *ApiController) ListWithTag(c *gin.Context) {
	api := &store.Api{}
	api.Id = c.Param("apiId")
	var err error

	apis, total, err := api.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	tag2Apis := make(map[string][]store.Api, 0)

	for _, api := range apis {
		if len(api.Tags) == 0 {
			tag := "default"
			oldTagApis := tag2Apis[tag]
			tag2Apis[tag] = append(oldTagApis, api)
		}

		for _, tag := range api.Tags {
			if tag == "" {
				tag = "default"
			}
			oldTagApis := tag2Apis[tag]
			tag2Apis[tag] = append(oldTagApis, api)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"total": total,
		"data":  tag2Apis,
	})
}

// @Summary List
// @Tags    Api
// @Success 200 {array} store.Api "api array"
// @Router /apis [get]
func (ctl *ApiController) List(c *gin.Context) {
	api := &store.Api{}
	api.Id = c.Param("apiId")
	var err error

	apis, total, err := api.List()
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
// @Success 200 {object} store.Api "api object"
// @Router /apis/{apiId} [get]
func (ctl *ApiController) Get(c *gin.Context) {
	api := &store.Api{}
	api.Id = c.Param("apiId")

	api, err := api.Get()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, api)
}
