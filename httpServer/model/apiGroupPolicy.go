//generate by gen
package model

import (
	"ioa/httpServer/app"
)

type ApiGroupPolicy struct {
	Common
	ApiGroupId string `json:"apiGroupId"`
	PolicyId   string `json:"policyId"`
}

func (apiGroupPolicy *ApiGroupPolicy) Insert() error {
	err := app.DB.Create(apiGroupPolicy).Error

	return err
}

func (apiGroupPolicy *ApiGroupPolicy) Patch() error {
	err := app.DB.Model(apiGroupPolicy).Updates(apiGroupPolicy).Error

	return err
}

func (apiGroupPolicy *ApiGroupPolicy) Update() error {
	err := app.DB.Save(apiGroupPolicy).Error

	return err
}

func (apiGroupPolicy *ApiGroupPolicy) Delete() error {
	return app.DB.Delete(apiGroupPolicy).Error
}

func (apiGroupPolicy *ApiGroupPolicy) List(rawQuery string, rawOrder string, offset int, limit int) (*[]ApiGroupPolicy, int, error) {
	apiGroupPolicys := []ApiGroupPolicy{}
	total := 0

	db := app.DB.Model(apiGroupPolicy)

	db, err := buildWhere(rawQuery, db)
	if err != nil {
		return &apiGroupPolicys, total, err
	}

	db, err = buildOrder(rawOrder, db)
	if err != nil {
		return &apiGroupPolicys, total, err
	}

	db.Offset(offset).
		Limit(limit).
		Find(&apiGroupPolicys).
		Count(&total)

	err = db.Error

	return &apiGroupPolicys, total, err
}

func (apiGroupPolicy *ApiGroupPolicy) Get() (*ApiGroupPolicy, error) {
	err := app.DB.Find(&apiGroupPolicy).Error

	return apiGroupPolicy, err
}
