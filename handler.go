package goma

import (
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/net/context"
)

func handleVersion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(Version))
}

// NewRouter creates gorilla/mux *Router for REST API.
func NewRouter(ctx context.Context) *mux.Router {
	r := mux.NewRouter()
	r.Path("/list").
		Name("list").
		Methods(http.MethodGet).
		HandlerFunc(handleList)

	r.Path("/register").
		Name("register").
		Methods(http.MethodPost).
		HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handleRegister(ctx, w, r)
		})

	r.Path("/monitor/{id:[0-9]+}").
		Name("monitor").
		HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handleMonitor(ctx, w, r)
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
