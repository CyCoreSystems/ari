package bridgemon

import (
	"sync"

	"github.com/CyCoreSystems/ari/v6"
)

// Monitor is a bridge monitor, which maintains bridge data.  It monitors an ARI bridge for events and keeps an internal cache of the bridge's data.
type Monitor struct {
	h *ari.BridgeHandle

	br *ari.BridgeData

	sub    ari.Subscription
	closed bool

	watchers  []chan *ari.BridgeData
	watcherMu sync.Mutex

	mu sync.Mutex
}

// New returns a new bridge monitor
func New(h *ari.BridgeHandle) *Monitor {
	sub := h.Subscribe(ari.Events.BridgeDestroyed, ari.Events.ChannelEnteredBridge, ari.Events.ChannelLeftBridge)

	m := &Monitor{
		h:   h,
		sub: sub,
	}

	// Monitor bridge events to keep data in sync
	go m.monitor()

	// Attempt to load initial bridge data; this may fail if the bridge has only
	// been staged, so ignore errors here
	data, _ := h.Data() // nolint
	m.updateData(data)

	return m
}

func (m *Monitor) monitor() {
	defer m.Close()

	for v := range m.sub.Events() {
		if v == nil {
			continue
		}

		switch v.GetType() {
		case ari.Events.BridgeDestroyed:
			e, ok := v.(*ari.BridgeDestroyed)
			if !ok {
				continue
			}

			m.updateData(&e.Bridge)

			return // bridge is destroyed; there will be no more events
		case ari.Events.ChannelEnteredBridge:
			e, ok := v.(*ari.ChannelEnteredBridge)
			if !ok {
				continue
			}

			m.updateData(&e.Bridge)
		case ari.Events.ChannelLeftBridge:
			e, ok := v.(*ari.ChannelLeftBridge)
			if !ok {
				continue
			}

			m.updateData(&e.Bridge)
		}
	}
}

func (m *Monitor) updateData(data *ari.BridgeData) {
	if data == nil {
		return
	}

	// Populate the bridge key in the bridge data, since Asterisk does not populate this field.
	if data.Key == nil {
		data.Key = m.h.Key()
	}

	// Update the stored data
	m.mu.Lock()
	m.br = data
	m.mu.Unlock()

	// Distribute new data to any watchers
	m.watcherMu.Lock()
	for _, w := range m.watchers {
		select {
		case w <- data:
		default:
		}
	}
	m.watcherMu.Unlock()
}

// Data returns the current bridge data
func (m *Monitor) Data() *ari.BridgeData {
	if m == nil {
		return nil
	}

	return m.br
}

// Handle returns the BridgeHandle which was used to create the bridge Monitor.
func (m *Monitor) Handle() *ari.BridgeHandle {
	if m == nil {
		return nil
	}

	return m.h
}

// Key returns the key of the monitored bridge
func (m *Monitor) Key() *ari.Key {
	if m == nil || m.h == nil {
		return nil
	}

	return m.h.Key()
}

// Watch returns a channel on which bridge data will be returned when events
// occur.  This channel will be closed when the bridge or the monitor is
// destoyed.
//
// NOTE:  the user should NEVER close this channel directly.
func (m *Monitor) Watch() <-chan *ari.BridgeData {
	ch := make(chan *ari.BridgeData)

	m.mu.Lock()
	defer m.mu.Unlock()

	if m.closed {
		close(ch)
		return ch
	}

	m.watcherMu.Lock()
	m.watchers = append(m.watchers, ch)
	m.watcherMu.Unlock()

	return ch
}

// Close shuts down a bridge monitor
func (m *Monitor) Close() {
	if m == nil {
		return
	}

	m.mu.Lock()
	if !m.closed {
		m.closed = true
		if m.sub != nil {
			m.sub.Cancel()
		}
	}
	m.mu.Unlock()

	m.watcherMu.Lock()

	for _, w := range m.watchers {
		close(w)
	}

	m.watchers = nil
	m.watcherMu.Unlock()
}
