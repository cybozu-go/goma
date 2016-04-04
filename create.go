package goma

import (
	"errors"
	"fmt"
	"time"

	"github.com/cybozu-go/goma/actions"
	"github.com/cybozu-go/goma/filters"
	"github.com/cybozu-go/goma/monitor"
	"github.com/cybozu-go/goma/probes"
)

const (
	typeKey         = "type"
	defaultInterval = 60 * time.Second
	defaultTimeout  = 59 * time.Second
)

// Errors for goma.
var (
	ErrBadName      = errors.New("bad monitor name")
	ErrNoType       = errors.New("no type")
	ErrInvalidType  = errors.New("invalid type")
	ErrInvalidRange = errors.New("invalid min/max range")
	ErrNoKey        = errors.New("no key")
)

// MonitorDefinition is a struct to load monitor definitions.
// TOML and JSON can be used.
type MonitorDefinition struct {
	Name     string                   `toml:"name" json:"name"`
	Probe    map[string]interface{}   `toml:"probe" json:"probe"`
	Filter   map[string]interface{}   `toml:"filter" json:"filter,omitempty"`
	Actions  []map[string]interface{} `toml:"actions" json:"actions"`
	Interval int                      `toml:"interval" json:"interval,omitempty"`
	Timeout  int                      `toml:"timeout" json:"timeout,omitempty"`
	Min      float64                  `toml:"min" json:"min,omitempty"`
	Max      float64                  `toml:"max" json:"max,omitempty"`
}

func getType(m map[string]interface{}) (t string, err error) {
	v, ok := m[typeKey]
	if !ok {
		err = ErrNoType
		return
	}
	s, ok := v.(string)
	if !ok {
		err = ErrInvalidType
		return
	}
	t = s
	return
}

func getParams(m map[string]interface{}) map[string]interface{} {
	nm := make(map[string]interface{})
	for k, v := range m {
		if k == typeKey {
			continue
		}
		nm[k] = v
	}
	return nm
}

// CreateMonitor creates a monitor from MonitorDefinition.
func CreateMonitor(d *MonitorDefinition) (*monitor.Monitor, error) {
	if len(d.Name) == 0 {
		return nil, ErrBadName
	}

	t, err := getType(d.Probe)
	if err != nil {
		return nil, err
	}
	probe, err := probes.Construct(t, getParams(d.Probe))
	if err != nil {
		return nil, fmt.Errorf("%s: %v in probe", d.Name, err)
	}

	var filter filters.Filter
	if d.Filter != nil {
		t, err = getType(d.Filter)
		if err != nil {
			return nil, err
		}
		f, err := filters.Construct(t, getParams(d.Filter))
		if err != nil {
			return nil, fmt.Errorf("%s: %v in filter", d.Name, err)
		}
		filter = f
	}

	var actors []actions.Actor
	for _, ad := range d.Actions {
		t, err = getType(ad)
		if err != nil {
			return nil, err
		}
		a, err := actions.Construct(t, getParams(ad))
		if err != nil {
			return nil, fmt.Errorf("%s: %v in action %s", d.Name, err, t)
		}
		actors = append(actors, a)
	}

	interval := time.Duration(d.Interval) * time.Second
	if interval == 0 {
		interval = defaultInterval
	}

	timeout := time.Duration(d.Timeout) * time.Second
	if timeout == 0 {
		timeout = defaultTimeout
	}

	if d.Min > d.Max {
		return nil, ErrInvalidRange
	}

	return monitor.NewMonitor(d.Name, probe, filter, actors,
		interval, timeout, d.Min, d.Max), nil
}
