//generate by gen
package model

import (
	"ioa/httpServer/app"
)

type Plugin struct {
	Common
	Name     string `json:"name"`
	Describe string `json:"describe"`
	FilePath string `json:"file_path"`
	Status   uint   `json:"status"`
}

func (plugin *Plugin) Insert() error {
	err := app.DB.Create(plugin).Error

	return err
}

func (plugin *Plugin) Patch() error {
	err := app.DB.Model(plugin).Updates(plugin).Error

	return err
}

func (plugin *Plugin) Update() error {
	err := app.DB.Save(plugin).Error

	return err
}

func (plugin *Plugin) Delete() error {
	return app.DB.Delete(plugin).Error
}

func (plugin *Plugin) List(rawQuery string, rawOrder string, offset int, limit int) (*[]Plugin, int, error) {
	plugins := []Plugin{}
	total := 0

	db := app.DB.Model(plugin)

	db, err := buildWhere(rawQuery, db)
	if err != nil {
		return &plugins, total, err
	}

	db, err = buildOrder(rawOrder, db)
	if err != nil {
		return &plugins, total, err
	}

	db.Offset(offset).
		Limit(limit).
		Find(&plugins).
		Count(&total)

	err = db.Error

	return &plugins, total, err
}

func (plugin *Plugin) Get() (*Plugin, error) {
	err := app.DB.Find(&plugin).Error

	return plugin, err
}
