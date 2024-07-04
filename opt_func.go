package request_locker

type MutexOpt func() Sync
type IdleOpt func() int64

func DefaultMutexOpt(mu Sync) func() Sync {
	return func() Sync {
		return mu
	}
}

func DefaultIdleOpt(t int64) func() int64 {
	return func() int64 {
		return t
	}
}
