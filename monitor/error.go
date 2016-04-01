package monitor

import "errors"

// Errors for monitors.
var (
	ErrRegistered    = errors.New("monitor has already been registered")
	ErrNotRegistered = errors.New("monitor has not been registered")
	ErrStarted       = errors.New("monitor has already been started")
)
