package router

import "net/http"

// Mux is...
type Mux struct {
	routes          map[string]*node
	NotFoundHandler http.HandlerFunc
	PanicHandler    func(http.ResponseWriter, *http.Request, interface{})
}

// Router is...
type Router struct {
	mux         *Mux
	prefix      string
	middlewares []func(http.Handler) http.Handler
}

// New is...
func New() *Router {
	return &Router{
		mux: &Mux{
			routes: make(map[string]*node),
		},
		middlewares: make([]func(http.Handler) http.Handler, 0),
	}
}

// Get is...
func (r *Router) Get(path interface{}, handle http.HandlerFunc) {
	r.Handle("GET", path, handle)
}

// Post is...
func (r *Router) Post(path interface{}, handle http.HandlerFunc) {
	r.Handle("POST", path, handle)
}

// Put is...
func (r *Router) Put(path interface{}, handle http.HandlerFunc) {
	r.Handle("PUT", path, handle)
}

// Patch is...
func (r *Router) Patch(path interface{}, handle http.HandlerFunc) {
	r.Handle("PATCH", path, handle)
}

// Delete is...
func (r *Router) Delete(path interface{}, handle http.HandlerFunc) {
	r.Handle("DELETE", path, handle)
}

// Options is...
func (r *Router) Options(path interface{}, handle http.HandlerFunc) {
	r.Handle("OPTIONS", path, handle)
}

// Handle is...
func (r *Router) Handle(method string, path interface{}, handle http.HandlerFunc) {
	p := path.(string)
	if p[0] != '/' {
		panic("can't handle relative path")
	}

	if r.mux.routes[method] == nil {
		r.mux.routes[method] = newNode(0)
	}

	r.mux.routes[method].insert(r.prefix+p, r.applyTo(handle))
}

// Subrouter is...
func (r *Router) Subrouter() *Router {
	subrouter := &Router{
		mux:         r.mux,
		middlewares: make([]func(http.Handler) http.Handler, len(r.middlewares)),
	}
	copy(subrouter.middlewares, r.middlewares)
	return subrouter
}

// Use is...
func (r *Router) Use(handler func(http.Handler) http.Handler) {
	r.middlewares = append(r.middlewares, handler)
}

func (r *Router) applyTo(handler http.Handler) http.Handler {
	for i := len(r.middlewares) - 1; i >= 0; i-- {
		handler = r.middlewares[i](handler)
	}
	return handler
}

func (m *Mux) recv(w http.ResponseWriter, req *http.Request) {
	if rcv := recover(); rcv != nil {
		m.PanicHandler(w, req, rcv)
	}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	m := r.mux
	if m.PanicHandler != nil {
		defer m.recv(w, req)
	}

	if m.routes[req.Method] != nil {
		path := req.URL.Path
		handle, found, params := m.routes[req.Method].find(path)
		if found {
			ctx := setContext(req)
			ctx.setParams(params)
			handle.ServeHTTP(w, req)
			return
		}
	}

	if m.NotFoundHandler != nil {
		m.NotFoundHandler(w, req)
	} else {
		http.NotFound(w, req)
	}
}
