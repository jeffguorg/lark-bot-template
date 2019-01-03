package web

import (
	"net/http"

	"github.com/go-chi/chi"
)

var router = func() *chi.Mux {
	r := chi.NewRouter()

	return r
}()

// Serve starts a http webserver
func Serve(addr string) error {
	return http.ListenAndServe(addr, router)
}
