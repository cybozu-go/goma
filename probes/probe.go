// Package probes provides API to implement goma probes.
package probes

import (
	"errors"
	"sync"

	"golang.org/x/net/context"
)

// Prober is the interface for probes.
type Prober interface {
	// Probe implements a probing method.
	//
	// The returned float64 value will be interpreted by the monitor
	// who run the probe.  If the probe encounters some error, it may
	// return non-nil error as well.
	//
	// ctx.Deadline() is always set.
	// Probe must return immediately when the ctx.Done() is closed.
	// Return values will not be used in such cases.
	Probe(ctx context.Context) (float64, error)

	// String returns a descriptive string for this probe.
	String() string
}

// Constructor is a function to create a probe.
//
// params are configuration options for the probe.
type Constructor func(params map[string]interface{}) (Prober, error)

// Errors for probes.
var (
	ErrNotFound = errors.New("probe not found")
)

var (
	registryLock = new(sync.Mutex)
	registry     = make(map[string]Constructor)
)

// Register registers a constructor of a kind of probes.
func Register(name string, ctor Constructor) {
	registryLock.Lock()
	defer registryLock.Unlock()

	if _, ok := registry[name]; ok {
		panic("duplicate probe entry: " + name)
	}

	registry[name] = ctor
}

// Construct constructs a named probe.
// This function is used internally in goma.
func Construct(name string, params map[string]interface{}) (Prober, error) {
	registryLock.Lock()
	ctor, ok := registry[name]
	registryLock.Unlock()

	if !ok {
		return nil, ErrNotFound
	}

	return ctor(params)
}
