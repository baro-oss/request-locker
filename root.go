package request_locker

import (
	"sync"
	"time"
)

type RootLocker struct {
	lockers map[interface{}]*Locker
	mu      sync.Locker
	config  *RootConfig
}

type RootConfig struct {
	idle   int64
	ticker time.Duration
}

func NewRootLocker(config *RootConfig) *RootLocker {
	return &RootLocker{
		lockers: make(map[interface{}]*Locker),
		mu:      &sync.Mutex{},
		config:  config,
	}
}

func (r *RootLocker) AddHolder(id interface{}, holder *SyncChannel[bool]) (err error) {
	r.mu.Lock()
	locker, ok := r.lockers[id]
	if !ok {
		locker = NewLocker(id, r.config.idle, r.config.ticker)
		locker.root = r
		r.lockers[id] = locker
		go locker.StartObserver()
	}
	r.mu.Unlock()
	locker.Assign(holder)
	return
}
