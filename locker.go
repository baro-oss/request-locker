package request_locker

import (
	"sync"
	"time"
)

type Locker struct {
	id          interface{}
	holder      *SyncChannel[bool]
	mu          sync.Locker
	lastChanged int64
	asIdleAt    int64
	ticker      time.Duration
	root        *RootLocker
}

func NewLocker(id interface{}, asIdleAt int64, ticker time.Duration) *Locker {
	locker := &Locker{
		id:       id,
		mu:       &sync.Mutex{},
		asIdleAt: asIdleAt,
		ticker:   ticker,
	}
	if asIdleAt == 0 {
		locker.asIdleAt = DefaultTimeLockerAsIdle
	}
	return locker
}

func (l *Locker) Assign(holder *SyncChannel[bool]) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.holder != nil {
		l.holder.Write(false)
	}
	l.holder = holder
	l.lastChanged = time.Now().UnixMilli()
}

func (l *Locker) Notify() {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.holder == nil {
		return
	}
	l.holder.Write(true)
	l.holder = nil
}

func (l *Locker) Close() {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.root == nil {
		return
	}
	delete(l.root.lockers, l.id)
	l.mu.Unlock()
}

func (l *Locker) StartObserver() {
	ticker := time.NewTicker(l.ticker)
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
