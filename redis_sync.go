package request_locker

import (
	"context"
	"time"

	"github.com/bsm/redislock"
)

type RedisSync struct {
	client  *redislock.Client
	lock    *redislock.Lock
	opts    *redislock.Options
	timeout time.Duration
	ttl     time.Duration
	key     string
}

func (r *RedisSync) Lock() {
	ctx, _ := context.WithTimeout(context.Background(), r.timeout)
	lock, err := r.client.Obtain(ctx, r.key, r.ttl, r.opts)
	if err != nil {
		return
	}
	r.lock = lock
}

func (r *RedisSync) Unlock() {
	if r.lock == nil {
		return
	}
	ctx, _ := context.WithTimeout(context.Background(), r.timeout)
	r.lock.Release(ctx)
}

func (r *RedisSync) TryLock() bool {
	return true
}
