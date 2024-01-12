package goma

import (
	"io"
	"net/http"
	"strings"

	"github.com/cybozu-go/log"
)

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
	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error("handleSetVerbosity", map[string]interface{}{
			"error": err.Error(),
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
	log.Info("new verbosity", map[string]interface{}{
		"level": level,
	})
}
