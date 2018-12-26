package main

import (
	"ioa"
	"ioa/httpServer"
	"ioa/httpServer/app"
	"ioa/httpServer/migrate"
)

func main() {
	app.InitConfig()
	app.InitDB()
	migrate.CreateTable()
	Ioa := ioa.New()
	go httpServer.Run(Ioa)

	Ioa.StartServer()
}
