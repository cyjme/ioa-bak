//generate by gen
package model

import (
	"ioa/httpServer/app"
)

type Target struct {
	Common
	ApiId  string `json:"apiId"`
	Scheme string `json:"scheme"`
	Method string `json:"method"`
	Host   string `json:"host"`
	Port   string `json:"port"`
	Path   string `json:"path"`
}

func (target *Target) Insert() error {
	err := app.DB.Create(target).Error

	return err
}

func (target *Target) Patch() error {
	err := app.DB.Model(target).Updates(target).Error

	return err
}

func (target *Target) Update() error {
	err := app.DB.Save(target).Error

	return err
}

func (target *Target) Delete() error {
	return app.DB.Delete(target).Error
}

func (target *Target) List(rawQuery string, rawOrder string, offset int, limit int) (*[]Target, int, error) {
	targets := []Target{}
	total := 0

	db := app.DB.Model(target)

	db, err := buildWhere(rawQuery, db)
	if err != nil {
		return &targets, total, err
	}

	db, err = buildOrder(rawOrder, db)
	if err != nil {
		return &targets, total, err
	}

	db.Offset(offset).
		Limit(limit).
		Find(&targets).
		Count(&total)

	err = db.Error

	return &targets, total, err
}

func (target *Target) Get() (*Target, error) {
	err := app.DB.Find(&target).Error

	return target, err
}
