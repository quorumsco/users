package gojimux

import (
	"net/http"

	"github.com/quorumsco/logs"
	"github.com/quorumsco/router"
	"github.com/zenazn/goji/web"
)

type Gojimux struct {
	Mux *web.Mux
}

func New() *Gojimux {
	var app = new(Gojimux)
	app.Mux = web.New()
	app.Mux.Use(InitContext)
	return app
}

func (app Gojimux) Get(path interface{}, handle http.HandlerFunc) {
	app.Mux.Get(path, app.putContext(handle))
}

func (app Gojimux) Post(path interface{}, handle http.HandlerFunc) {
	app.Mux.Post(path, app.putContext(handle))
}

func (app Gojimux) Put(path interface{}, handle http.HandlerFunc) {
	app.Mux.Put(path, handle)
}

func (app Gojimux) Patch(path interface{}, handle http.HandlerFunc) {
	app.Mux.Patch(path, app.putContext(handle))
}

func (app Gojimux) Delete(path interface{}, handle http.HandlerFunc) {
	app.Mux.Delete(path, app.putContext(handle))
}

func (app Gojimux) Options(path interface{}, handle http.HandlerFunc) {
	app.Mux.Options(path, app.putContext(handle))
}

func (app Gojimux) Use(handler func(http.Handler) http.Handler) {
	app.Mux.Use(handler)
}

func (app Gojimux) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	app.Mux.ServeHTTP(w, req)
}

func (app Gojimux) Serve(listen string) error {
	logs.Info("listening on http://%s", listen)
	return http.ListenAndServe(listen, app.Mux)
}

func (app Gojimux) putContext(handle http.HandlerFunc) func(web.C, http.ResponseWriter, *http.Request) {
	fn := func(c web.C, w http.ResponseWriter, r *http.Request) {
		param := new(router.Param)
		for name, value := range c.URLParams {
			param.Name = name
			param.Value = value
			router.Context(r).Params = append(router.Context(r).Params, *param)
		}
		handle(w, r)
	}
	return fn
}
