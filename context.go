package ioa

import "net/http"

type Context struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request
	Response       *http.Response
	Api            *Api
	Next           bool
}

func (c *Context) Cancel() {
	c.Next = false
}
