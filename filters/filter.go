// Package filters provides API to implement goma filters.
package filters

import (
	"errors"
	"sync"
)

// Filter is the interface for filters.
type Filter interface {
	// Init is called when goma starts monitoring.
	Init()

	// Put receives a return value from a probe, and returns a filtered value.
	Put(f float64) float64

	// String returns a descriptive string for this filter.
	String() string
}

// Constructor is a function to create a filter.
//
// params are configuration options for the probe.
type Constructor func(params map[string]interface{}) (Filter, error)

// Errors for filters.
var (
	ErrNotFound = errors.New("filter not found")
)

var (
	registryLock = new(sync.Mutex)
	registry     = make(map[string]Constructor)
)

// Register registers a constructor of a kind of filters.
func Register(name string, ctor Constructor) {
	registryLock.Lock()
	defer registryLock.Unlock()

	if _, ok := registry[name]; ok {
		panic("duplicate filter entry: " + name)
	}

	registry[name] = ctor
}

// Construct constructs a named filter.
// This function is used internally in goma.
func Construct(name string, params map[string]interface{}) (Filter, error) {
	registryLock.Lock()
	ctor, ok := registry[name]
	registryLock.Unlock()

	if !ok {
		return nil, ErrNotFound
	}

	return ctor(params)
}
