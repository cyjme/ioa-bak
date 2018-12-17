package main

import (
	"ioa/count"
	"log"
	"net/http"
	"sync/atomic"
)

func Run(w http.ResponseWriter, r *http.Request) {
	atomic.AddInt64(&count.Count, 1)
	log.Println("插件开始处理")
	log.Println("current count", count.Count)
}
