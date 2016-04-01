package goma

import (
	"encoding/json"
	"fmt"
	"mime"
	"net/http"

	"github.com/cybozu-go/goma/monitor"
	"github.com/cybozu-go/log"
	"golang.org/x/net/context"
)

func handleRegister(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	mt, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if mt != "application/json" {
		http.Error(w, "bad content type", http.StatusBadRequest)
		return
	}

	d := json.NewDecoder(r.Body)
	var md MonitorDefinition
	if err := d.Decode(&md); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	m, err := CreateMonitor(&md)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// ignoring error is safe here.
	monitor.Register(m)
	log.Info("new monitor", map[string]interface{}{
		"_id":   m.ID(),
		"_name": m.Name(),
	})
	m.Start(ctx)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(fmt.Sprintf("%d", m.ID())))
}
