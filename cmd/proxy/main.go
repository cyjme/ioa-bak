package main

import (
	"ioa"
	"ioa/store"
)

func main() {
	store.Init()
	config := ioa.ReadConfig()
	Ioa := ioa.New(config)
	Ioa.StartServer()
}
