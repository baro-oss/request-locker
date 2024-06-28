package request_locker

import (
	"fmt"
	"time"
)

type Locker struct {
	id          interface{}
	holder      chan bool
	mu          Sync
	lastChanged int64
	asIdleAt    int64
	root        *Root
}

type LockerOptFunc func(*Locker)

func NewLocker(id interface{}, opts ...LockerOptFunc) (*Locker, error) {
	l := &Locker{id: id}
	for _, opt := range opts {
		opt(l)
	}
	if l.mu == nil {
		return nil, ErrInvalidSyncLocker
	}
	if l.asIdleAt == 0 {
		l.asIdleAt = DefaultTimeLockerAsIdle
	}
	return l, nil
}

func (l *Locker) Assign(holder chan bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.holder != nil {
		l.holder <- false
	}
	l.holder = holder
	l.lastChanged = time.Now().UnixMilli()
}

func (l *Locker) Notify() {
	defer func() {
		// for case timeout => send to a closed channel
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.holder == nil {
		return
	}
	l.holder <- true
	l.holder = nil
}

func (l *Locker) Close() {
	l.mu.Lock()
	delete(l.root.lockers, l.id)
	l.mu.Unlock()
}

func (l *Locker) StartObserver() {
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case t := <-ticker.C:
			if t.UnixMilli()-l.lastChanged < l.asIdleAt {
				continue
			}
			l.Close()
			l.Notify()
			return
		}
	}
}
