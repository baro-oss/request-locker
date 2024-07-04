package request_locker

import "sync"

type Sync interface {
	Lock()
	Unlock()
	TryLock() bool
}

type SyncChannel[T any] struct {
	ch       chan T
	isClosed bool
	mu       sync.Locker
}

func NewSyncChannel[T any](innerChan chan T) *SyncChannel[T] {
	syncChannel := &SyncChannel[T]{
		isClosed: false,
		mu:       &sync.Mutex{},
		ch:       innerChan,
	}
	if _, ok := <-innerChan; !ok {
		syncChannel.isClosed = true
	}
	return syncChannel
}

func (c *SyncChannel[T]) IsClosed() bool {
	return c.isClosed
}

func (c *SyncChannel[T]) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.isClosed {
		return ErrorCloseClosedChannel
	}
	close(c.ch)
	return nil
}

func (c *SyncChannel[T]) Write(value T) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.isClosed {
		return ErrorWriteClosedChannel
	}
	c.ch <- value
	return nil
}

func (c *SyncChannel[T]) Read() (T, error) {
	var t T
	if c.isClosed {
		return t, ErrorReadClosedChannel
	}
	t, ok := <-c.ch
	c.isClosed = !ok
	return t, nil
}
