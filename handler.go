package goma

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/cybozu-go/log"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
)

func handleVersion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(Version))
}

func handleVerbosity(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte(log.LevelName(log.DefaultLogger().Threshold())))
		return
	}

	if r.Method != http.MethodPut && r.Method != http.MethodPost {
		http.Error(w, "Bad method", http.StatusBadRequest)
		return
	}

	// for PUT or POST, set new verbosity.
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error("handleSetVerbosity", map[string]interface{}{
			"_err": err.Error(),
		})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	level := strings.TrimSpace(string(data))
	err = log.DefaultLogger().SetThresholdByName(level)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

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

	r.Path("/version").
		Name("version").
		Methods(http.MethodGet).
		HandlerFunc(handleVersion)

	r.Path("/verbosity").
		Name("verbosity").
		HandlerFunc(handleVerbosity)
	return r
}
