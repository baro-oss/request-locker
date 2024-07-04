package request_locker

import (
	"time"
)

type Locker struct {
	id          interface{}
	holder      *SyncChannel[bool]
	mu          Sync
	lastChanged int64
	asIdleAt    int64
	root        *RootLocker
}

func NewLocker(id interface{}, asIdleAt int64, mu Sync) *Locker {
	return &Locker{
		id:       id,
		mu:       mu,
		asIdleAt: asIdleAt,
	}
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
