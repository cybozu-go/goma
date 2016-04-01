package goma

import (
	_log "log"
	"net"
	"net/http"
	"time"

	"github.com/cybozu-go/log"
	"github.com/facebookgo/httpdown"
	"golang.org/x/net/context"
)

const (
	defaultReadTimeout  = 10 * time.Second
	defaultWriteTimeout = 10 * time.Second

	// Version may be used for REST API version checks in future.
	Version = "0.1"

	// VersionHeader is the HTTP request header for Version.
	VersionHeader = "X-Goma-Version"
)

// Serve runs REST API server until ctx.Done() is closed.
func Serve(ctx context.Context, l net.Listener) error {
	hd := httpdown.HTTP{}
	logger := _log.New(log.DefaultLogger().Writer(log.LvError), "[http]", 0)
	s := &http.Server{
		Handler:      NewRouter(ctx),
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
		ErrorLog:     logger,
	}
	hs := hd.Serve(s, l)

	waiterr := make(chan error, 1)
	go func() {
		defer close(waiterr)
		waiterr <- hs.Wait()
	}()

	select {
	case err := <-waiterr:
		return err

	case <-ctx.Done():
		if err := hs.Stop(); err != nil {
			return err
		}
		return <-waiterr
	}
}
