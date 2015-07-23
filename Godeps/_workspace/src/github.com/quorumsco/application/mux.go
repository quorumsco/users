package application

import "net/http"

// Mux is the interface requests multiplexers must implement to be used
// with Application.
type Mux interface {
	Get(path interface{}, handle http.HandlerFunc)
	Post(path interface{}, handle http.HandlerFunc)
	Put(path interface{}, handle http.HandlerFunc)
	Patch(path interface{}, handle http.HandlerFunc)
	Delete(path interface{}, handle http.HandlerFunc)
	Options(path interface{}, handle http.HandlerFunc)
	Use(handler func(http.Handler) http.Handler)
	ServeHTTP(w http.ResponseWriter, req *http.Request)
}
