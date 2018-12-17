package migrate

import (
	"ioa/httpServer/app"
	"ioa/httpServer/model"
)

func CreateTable() {
	app.DB.AutoMigrate(
		&model.Api{},
		&model.ApiGroup{},
		&model.ApiPolicy{},
		&model.ApiGroupPolicy{},
		&model.Param{},
		&model.Policy{},
		&model.Target{},
		&model.Plugin{},
		//!!do not delete the line, gen generate code at here

	)
}
