//generate by gen
package model

import (
	"httpServer/app"
)

type Policy struct {
	Common
	Name             string `json:"name"`
	Describe         string `json:"describe"`
	Plugin_file_name string `json:"plugin_file_name"`
}

func (policy *Policy) Insert() error {
	err := app.DB.Create(policy).Error

	return err
}

func (policy *Policy) Patch() error {
	err := app.DB.Model(policy).Updates(policy).Error

	return err
}

func (policy *Policy) Update() error {
	err := app.DB.Save(policy).Error

	return err
}

func (policy *Policy) Delete() error {
	return app.DB.Delete(policy).Error
}

func (policy *Policy) List(rawQuery string, rawOrder string, offset int, limit int) (*[]Policy, int, error) {
	policys := []Policy{}
	total := 0

	db := app.DB.Model(policy)

	db, err := buildWhere(rawQuery, db)
	if err != nil {
		return &policys, total, err
	}

	db, err = buildOrder(rawOrder, db)
	if err != nil {
		return &policys, total, err
	}

	db.Offset(offset).
		Limit(limit).
		Find(&policys).
		Count(&total)

	err = db.Error

	return &policys, total, err
}

func (policy *Policy) Get() (*Policy, error) {
	err := app.DB.Find(&policy).Error

	return policy, err
}
