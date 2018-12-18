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
	go httpServer.Run()
	ioa.StartServer()
}
