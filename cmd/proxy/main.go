package main

import (
	"flag"
	"ioa"
	"ioa/store"
)

func main() {
	config := ioa.ReadConfig()

	flag.StringVar(&config.Proxy.Host, "host", "0.0.0.0", "")
	flag.StringVar(&config.Proxy.Port, "port", "9992", "")
	flag.IntVar(&config.Proxy.MaxIdleConns, "maxIdleConns", 10000, "")
	flag.IntVar(&config.Proxy.MaxIdleConnsPerHost, "maxIdleConnsPerHost", 10000, "")
	flag.Parse()

	store.Init()
	Ioa := ioa.New(config)
	Ioa.StartServer()
}
