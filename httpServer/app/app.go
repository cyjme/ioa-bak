package app

import "golang.org/x/time/rate"

var Limiter = rate.NewLimiter(1, 5)

var Count int