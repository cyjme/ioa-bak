package main

import (
	"ioa"
	"ioa/httpServer/app"
	"ioa/store"
)

func main() {
	store.Init()
	app.InitConfig()
	Ioa := ioa.New()
	Ioa.StartServer()
}
