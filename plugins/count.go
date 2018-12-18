package main

import (
	"ioa"
	"log"
	"net/http"
	"sync/atomic"
)

func Run(w http.ResponseWriter, r *http.Request) {
	atomic.AddInt64(&ioa.Count, 1)
	log.Println("插件开始处理")
	log.Println("current count", ioa.Count)
}
