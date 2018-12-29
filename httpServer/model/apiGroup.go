//generate by gen
package model

import (
	"github.com/jinzhu/gorm"
	"ioa/httpServer/app"
)

type ApiGroup struct {
	Common
	Name      string `json:"name"`
	Describe  string `json:"describe"`
	IsDefault bool   `json:"isDefault"`

	Apis     []Api  `json:"apis"`
	Policies string `json:"policies"`
	Plugins  string `json:"plugins"`
}

func (apiGroup *ApiGroup) Insert() error {
	err := app.DB.Create(apiGroup).Error

	return err
}

func (apiGroup *ApiGroup) Patch() error {
	err := app.DB.Model(apiGroup).Updates(apiGroup).Error

	return err
}

func (apiGroup *ApiGroup) Update() error {
	err := app.DB.Save(apiGroup).Error

	return err
}

func (apiGroup *ApiGroup) Delete() error {
	return app.DB.Delete(apiGroup).Error
}

func (apiGroup *ApiGroup) List(rawQuery string, rawOrder string, offset int, limit int) (*[]ApiGroup, int, error) {
	apiGroups := []ApiGroup{}
	total := 0

	db := app.DB.Model(apiGroup).Preload("Apis",
		func(db *gorm.DB) *gorm.DB {
			return db.Order("api.created_at ASC")
		}).Preload("Apis.Targets",
		func(db *gorm.DB) *gorm.DB {
			return db.Order("target.created_at ASC")
		}).Preload("Apis.Params")

	db, err := buildWhere(rawQuery, db)
	if err != nil {
		return &apiGroups, total, err
	}

	db, err = buildOrder(rawOrder, db)
	if err != nil {
		return &apiGroups, total, err
	}

	db.Offset(offset).
		Limit(limit).
		Find(&apiGroups).
		Count(&total)

	err = db.Error

	return &apiGroups, total, err
}

func (apiGroup *ApiGroup) Get() (*ApiGroup, error) {
	err := app.DB.Model(apiGroup).Preload("Apis",
		func(db *gorm.DB) *gorm.DB {
			return db.Order("api.created_at ASC")
		}).Preload("Apis.Targets",
		func(db *gorm.DB) *gorm.DB {
			return db.Order("target.created_at ASC")
		}).Preload("Apis.Params").Find(&apiGroup).Error

	return apiGroup, err
}
