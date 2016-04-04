package monitor

import (
	"fmt"
	"sync"
	"time"

	"github.com/cybozu-go/goma/actions"
	"github.com/cybozu-go/goma/filters"
	"github.com/cybozu-go/goma/probes"
	"github.com/cybozu-go/log"
	"golang.org/x/net/context"
)

// Monitor is a unit of monitoring.
//
// It consists of a (configured) probe, zero or one filter, and one or
// more actions.  goma will invoke Prover.Probe periodically at given
// interval.
type Monitor struct {
	id       int
	name     string
	probe    probes.Prober
	filter   filters.Filter
	actors   []actions.Actor
	interval time.Duration
	timeout  time.Duration
	min      float64
	max      float64
	failedAt *time.Time

	// goroutine management
	lock   sync.Mutex
	cancel context.CancelFunc
	done   <-chan struct{}
}

// NewMonitor creates and initializes a monitor.
//
// name can be any descriptive string for the monitor.
// p and a should not be nil.  f may be nil.
// interval is the interval between probes.
// timeout is the maximum duration for a probe to run.
// min and max defines the range for normal probe results.
func NewMonitor(
	name string,
	p probes.Prober,
	f filters.Filter,
	a []actions.Actor,
	interval, timeout time.Duration,
	min, max float64) *Monitor {
	return &Monitor{
		id:       uninitializedID,
		name:     name,
		probe:    p,
		filter:   f,
		actors:   a,
		interval: interval,
		timeout:  timeout,
		min:      min,
		max:      max,
	}
}

// Start starts monitoring.
// If already started, this returns a non-nil error.
func (m *Monitor) Start(ctx context.Context) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	if m.done != nil {
		return ErrStarted
	}

	done := make(chan struct{})
	m.done = done
	ctx, cancel := context.WithCancel(ctx)
	m.cancel = cancel

	go run(ctx, m, done)

	log.Info("monitor started", map[string]interface{}{
		"_monitor": m.name,
	})

	return nil
}

// Stop stops monitoring.
func (m *Monitor) Stop() {
	m.lock.Lock()
	defer m.lock.Unlock()

	if m.done == nil {
		return
	}

	log.Debug("monitor is stopping", map[string]interface{}{
		"_monitor": m.name,
	})

	m.cancel()

	<-m.done
	m.done = nil

	m.failedAt = nil

	log.Info("monitor stopped", map[string]interface{}{
		"_monitor": m.name,
	})
}

func callProbe(ctx context.Context, p probes.Prober, timeout time.Duration) float64 {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return p.Probe(ctx)
}

func run(ctx context.Context, m *Monitor, done chan<- struct{}) {
	if m.filter != nil {
		m.filter.Init()
	}
	for _, a := range m.actors {
		a.Init(m.name)
	}

	for {
		// create a timer before starting probe.
		// This way, we can keep consistent interval between probes.
		t := time.After(m.interval)

		v := callProbe(ctx, m.probe, m.timeout)

		// check cancel
		select {
		case <-ctx.Done():
			done <- struct{}{}
			return
		default:
			// not canceled
		}

		if m.filter != nil {
			v = m.filter.Put(v)
		}

		if (v < m.min) || (m.max < v) {
			if m.failedAt == nil {
				now := time.Now()
				m.failedAt = &now
				for _, a := range m.actors {
					a.Fail(m.name, v)
				}
				log.Warn("monitor failure", map[string]interface{}{
					"_monitor": m.name,
					"_value":   fmt.Sprint(v),
				})
			}
		} else {
			if m.failedAt != nil {
				d := time.Since(*m.failedAt)
				for _, a := range m.actors {
					a.Recover(m.name, d)
				}
				m.failedAt = nil
				log.Warn("monitor recovery", map[string]interface{}{
					"_monitor":  m.name,
					"_duration": d.Seconds(),
				})
			}
		}

		select {
		case <-ctx.Done():
			done <- struct{}{}
			return
		case <-t:
			// interval timer expires
		}
	}
}

// ID returns the monitor ID.
//
// ID is valid only after registration.
func (m *Monitor) ID() int {
	return m.id
}

// Name returns the name of the monitor.
func (m *Monitor) Name() string {
	return m.name
}

// String is the same as Name.
func (m *Monitor) String() string {
	return m.name
}

// Failing returns true if the monitor is detecting a failure.
func (m *Monitor) Failing() bool {
	return m.failedAt != nil
}

// Running returns true if the monitor is running.
func (m *Monitor) Running() bool {
	m.lock.Lock()
	defer m.lock.Unlock()

	return m.done != nil
}
