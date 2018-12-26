package httpServer

import (
	"ioa"
	"ioa/httpServer/router"
)

func Run(ioa *ioa.Ioa) {
	router.Start(ioa)
}
