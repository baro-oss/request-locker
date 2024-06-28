package request_locker

import "errors"

const (
	DefaultTimeLockerAsIdle = 1000
)

var (
	ErrInvalidSyncLocker = errors.New("invalid sync locker")
)
