package goma

import (
	"net/http"

	"github.com/gorilla/mux"
)

func handleVersion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(Version))
}

// NewRouter creates gorilla/mux *Router for REST API.
func NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.Path("/list").
		Name("list").
		Methods(http.MethodGet).
		HandlerFunc(handleList)

	r.Path("/register").
		Name("register").
		Methods(http.MethodPost).
		HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handleRegister(w, r)
		})

	r.Path("/monitor/{id:[0-9]+}").
		Name("monitor").
		HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handleMonitor(w, r)
		})

	r.Path("/verbosity").
		Name("verbosity").
		HandlerFunc(handleVerbosity)

	r.Path("/version").
		Name("version").
		Methods(http.MethodGet).
		HandlerFunc(handleVersion)

	return r
}
