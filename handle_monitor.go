package goma

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/net/context"

	"github.com/cybozu-go/goma/monitor"
	"github.com/gorilla/mux"
)

// MonitorInfo represents status of a monitor.
// This is used by show and list commands.
type MonitorInfo struct {
	ID      int    `json:"id,string"`
	Name    string `json:"name"`
	Running bool   `json:"running"`
	Failing bool   `json:"failing"`
}

func handleMonitor(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	// guaranteed no error by mux.
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	m := monitor.FindMonitor(id)
	if m == nil {
		http.NotFound(w, r)
		return
	}

	if r.Method == http.MethodGet {
		mi := &MonitorInfo{
			ID:      m.ID(),
			Name:    m.Name(),
			Running: m.Running(),
			Failing: m.Failing(),
		}
		data, err := json.Marshal(mi)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write(data)
		return
	}

	if r.Method == http.MethodDelete {
		m.Stop()
		monitor.Unregister(m)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "invalid method", http.StatusBadRequest)
		return
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	switch strings.TrimSpace(string(data)) {
	case "start":
		if err := m.Start(ctx); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	case "stop":
		m.Stop()
	default:
		http.Error(w, "unknown action", http.StatusBadRequest)
	}
}
