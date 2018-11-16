package goma

import (
	"net/http"
	"time"

	"github.com/cybozu-go/well"
)

const (
	defaultReadTimeout  = 10 * time.Second
	defaultWriteTimeout = 10 * time.Second

	// Version may be used for REST API version checks in future.
	Version = "1.0"

	// VersionHeader is the HTTP request header for Version.
	VersionHeader = "X-Goma-Version"
)

// Serve runs REST API server until the global environment is canceled.
func Serve(addr string) error {
	s := &well.HTTPServer{
		Server: &http.Server{
			Addr:         addr,
			Handler:      NewRouter(),
			ReadTimeout:  defaultReadTimeout,
			WriteTimeout: defaultWriteTimeout,
		},
	}
	return s.ListenAndServe()
}
