package request_locker

import (
	"errors"
	"time"
)

const (
	DefaultTimeLockerAsIdle = 1000
	DefaultTicker           = 20 * time.Millisecond
	DefaultTTL              = 1 * time.Second
)

var (
	ErrInvalidSyncLocker = errors.New("invalid sync locker")

	ErrorReadClosedChannel  = errors.New("read from a closed channel")
	ErrorCloseClosedChannel = errors.New("close a closed channel")
	ErrorWriteClosedChannel = errors.New("write to a closed channel")
)
