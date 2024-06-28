package request_locker

import (
	"sync"
)

type Root struct {
	lockers map[interface{}]*Locker
	mu      sync.Locker
	config  *RootConfig
}

type RootConfig struct {
	LockerOptFuncs []LockerOptFunc
}

func NewRoot(config *RootConfig) *Root {
	return &Root{
		lockers: make(map[interface{}]*Locker),
		mu:      &sync.Mutex{},
		config:  config,
	}
}

func (r *Root) AddHolder(id interface{}, holder chan bool) (err error) {
	r.mu.Lock()
	locker, ok := r.lockers[id]
	if !ok {
		locker, err = NewLocker(id, r.config.LockerOptFuncs...)
		if err != nil {
			return
		}
		locker.root = r
		r.lockers[id] = locker
		go locker.StartObserver()
	}
	r.mu.Unlock()
	locker.Assign(holder)
	return
}
