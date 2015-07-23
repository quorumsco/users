package router

import (
	"io"
	"net/http"
)

// RequestContext is...
type RequestContext struct {
	Env    map[interface{}]interface{}
	Params []param
	req    *http.Request
}

type iogoContextReadCloser struct {
	io.ReadCloser
	context *RequestContext
}

func setContext(req *http.Request) *RequestContext {
	c, ok := req.Body.(iogoContextReadCloser)
	if !ok {
		c = iogoContextReadCloser{
			ReadCloser: req.Body,
			context:    &RequestContext{Env: make(map[interface{}]interface{})},
		}
		req.Body = c
	}
	return c.context
}

// Context is...
func Context(req *http.Request) *RequestContext {
	return req.Body.(iogoContextReadCloser).context
}

// Param is...
func (c *RequestContext) Param(name string) string {
	for _, e := range c.Params {
		if e.name == name {
			return e.value
		}
	}
	return ""
}

// SetParams is...
func (c *RequestContext) setParams(params []param) {
	c.Params = params
}
