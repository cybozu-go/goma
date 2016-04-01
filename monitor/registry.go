package monitor

import "sync"

const (
	uninitializedID = -1
)

var (
	registryLock  = new(sync.Mutex)
	registry      = make(map[int]*Monitor)
	registryIndex int
)

// Register registers a monitor.
func Register(m *Monitor) error {
	if m.id != uninitializedID {
		return ErrRegistered
	}

	registryLock.Lock()
	defer registryLock.Unlock()

	m.id = registryIndex
	registry[registryIndex] = m
	registryIndex++
	return nil
}

// FindMonitor looks up a monitor in the registry.
// If not found, nil is returned.
func FindMonitor(id int) *Monitor {
	registryLock.Lock()
	defer registryLock.Unlock()

	return registry[id]
}

// Unregister removes a monitor from the registry.
// The monitor should have stopped.
func Unregister(m *Monitor) error {
	if m.id == uninitializedID {
		return ErrNotRegistered
	}

	registryLock.Lock()
	defer registryLock.Unlock()

	delete(registry, m.id)
	m.id = uninitializedID
	return nil
}

// ListMonitors returns a list of monitors ordered by ID (ascending).
func ListMonitors() []*Monitor {
	registryLock.Lock()
	defer registryLock.Unlock()

	l := make([]*Monitor, 0, len(registry))

	for i := 0; i < registryIndex; i++ {
		if m, ok := registry[i]; ok {
			l = append(l, m)
		}
	}

	return l
}
