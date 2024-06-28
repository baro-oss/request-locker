package request_locker

type Sync interface {
	Lock()
	Unlock()
	TryLock() bool
}
