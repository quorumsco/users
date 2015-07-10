package application

import (
	"net/http"

	"github.com/iogo-framework/router"
)

// Apply is a middleware which make the app context available to the
// requests handlers through the router. Obviously, this stands for
// iogo-framework/router.
func (app *Application) Apply(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		router.Context(r).Env["Application"] = app
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

// App retrieves the Application stored by the middleware Apply.
func App(r *http.Request) *Application {
	return router.Context(r).Env["Application"].(*Application)
}
