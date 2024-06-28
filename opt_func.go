package request_locker

func LockerSyncOpt(sync Sync) func(*Locker) {
	return func(locker *Locker) {
		if rs, ok := sync.(*RedisSync); ok {
			rs.key = locker.id.(string)
			locker.mu = rs
			return
		}

		locker.mu = sync
	}
}

func LockerIdleTimeOpt(t int64) func(*Locker) {
	return func(locker *Locker) {
		locker.asIdleAt = t
	}
}
