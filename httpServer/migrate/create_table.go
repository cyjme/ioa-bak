package migrate

import (
	"httpServer/app"
	"httpServer/model"
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
	//!!do not delete the line, gen generate code at here

	)
}
