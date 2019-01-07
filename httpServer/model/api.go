//generate by gen
package model

import (
	"github.com/jinzhu/gorm"
	"ioa/httpServer/app"
)

type Api struct {
	Common
	ApiGroupId string `json:"apiGroupId"`
	Name       string `json:"name"`
	Describe   string `json:"describe"`
	Path       string `json:"path"`
	Method     string `json:"method"`
	Status     string `json:"status"`

	Targets      []Target `json:"targets"`
	Params       []Param  `json:"params"`
	Policies     string   `json:"policies"`
	Plugins      string   `json:"plugins"`
}

func (api *Api) Insert() error {
	err := app.DB.Create(api).Error

	return err
}

func (api *Api) Patch() error {
	err := app.DB.Model(api).Updates(api).Error

	return err
}

func (api *Api) Update() error {
	err := app.DB.Save(api).Error

	return err
}

func (api *Api) Delete() error {
	return app.DB.Delete(api).Error
}

func (api *Api) List(rawQuery string, rawOrder string, offset int, limit int) (*[]Api, int, error) {
	apis := []Api{}
	total := 0

	db := app.DB.Model(api).
		Preload("Targets",
			func(db *gorm.DB) *gorm.DB {
				return db.Order("target.created_at ASC")
			}).
		Preload("Params",
			func(db *gorm.DB) *gorm.DB {
				return db.Order("param.created_at ASC")
			})

	db, err := buildWhere(rawQuery, db)
	if err != nil {
		return &apis, total, err
	}

	db, err = buildOrder(rawOrder, db)
	if err != nil {
		return &apis, total, err
	}

	db.Offset(offset).
		Limit(limit).
		Find(&apis).
		Count(&total)

	err = db.Error

	return &apis, total, err
}

func (api *Api) Get() (*Api, error) {
	err := app.DB.Model(api).
		Preload("Targets",
			func(db *gorm.DB) *gorm.DB {
				return db.Order("target.created_at ASC")
			}).
		Preload("Params",
			func(db *gorm.DB) *gorm.DB {
				return db.Order("param.created_at ASC")
			}).
		Find(&api).Error

	return api, err
}
