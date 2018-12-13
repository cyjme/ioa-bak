//generate by gen
package model

import (
	"httpServer/app"
)

type ApiPolicy struct {
	Common
	ApiId    string `json:"apiId"`
	PolicyId string `json:"policyId"`
	Priority int    `json:"priority"`
}

func (apiPolicy *ApiPolicy) Insert() error {
	err := app.DB.Create(apiPolicy).Error

	return err
}

func (apiPolicy *ApiPolicy) Patch() error {
	err := app.DB.Model(apiPolicy).Updates(apiPolicy).Error

	return err
}

func (apiPolicy *ApiPolicy) Update() error {
	err := app.DB.Save(apiPolicy).Error

	return err
}

func (apiPolicy *ApiPolicy) Delete() error {
	return app.DB.Delete(apiPolicy).Error
}

func (apiPolicy *ApiPolicy) List(rawQuery string, rawOrder string, offset int, limit int) (*[]ApiPolicy, int, error) {
	apiPolicys := []ApiPolicy{}
	total := 0

	db := app.DB.Model(apiPolicy)

	db, err := buildWhere(rawQuery, db)
	if err != nil {
		return &apiPolicys, total, err
	}

	db, err = buildOrder(rawOrder, db)
	if err != nil {
		return &apiPolicys, total, err
	}

	db.Offset(offset).
		Limit(limit).
		Find(&apiPolicys).
		Count(&total)

	err = db.Error

	return &apiPolicys, total, err
}

func (apiPolicy *ApiPolicy) Get() (*ApiPolicy, error) {
	err := app.DB.Find(&apiPolicy).Error

	return apiPolicy, err
}
