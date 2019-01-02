package app

import "golang.org/x/time/rate"

var Limiter = rate.NewLimiter(2, 5)

var Count int