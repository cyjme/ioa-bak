//generate by gen
package model

import (
	"ioa/httpServer/app"
)

type Param struct {
	Common
	ApiId      string `json:"apiId"`
	Location   string `json:"location"`
	Name       string `json:"name"`
	TargetName string `json:"targetName"`
}

func (param *Param) Insert() error {
	err := app.DB.Create(param).Error

	return err
}

func (param *Param) Patch() error {
	err := app.DB.Model(param).Updates(param).Error

	return err
}

func (param *Param) Update() error {
	err := app.DB.Save(param).Error

	return err
}

func (param *Param) Delete() error {
	return app.DB.Delete(param).Error
}

func (param *Param) List(rawQuery string, rawOrder string, offset int, limit int) (*[]Param, int, error) {
	params := []Param{}
	total := 0

	db := app.DB.Model(param)

	db, err := buildWhere(rawQuery, db)
	if err != nil {
		return &params, total, err
	}

	db, err = buildOrder(rawOrder, db)
	if err != nil {
		return &params, total, err
	}

	db.Offset(offset).
		Limit(limit).
		Find(&params).
		Count(&total)

	err = db.Error

	return &params, total, err
}

func (param *Param) Get() (*Param, error) {
	err := app.DB.Find(&param).Error

	return param, err
}
