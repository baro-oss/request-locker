package request_locker

import (
	"context"
	"time"

	"github.com/bsm/redislock"
	"github.com/redis/go-redis/v9"
)

type DistributedLocker struct {
	redisClient    *redis.Client
	redisLock      *redislock.Client
	ttl            time.Duration
	tickerDuration time.Duration
}

func NewDistributedLocker(rsClient *redis.Client, rsLock *redislock.Client, ttl, ticketDuration time.Duration) *DistributedLocker {
	distributedLocker := &DistributedLocker{
		redisClient:    rsClient,
		redisLock:      rsLock,
		ttl:            ttl,
		tickerDuration: ticketDuration,
	}
	if distributedLocker.ttl == 0 {
		distributedLocker.ttl = DefaultTTL
	}
	if distributedLocker.tickerDuration == 0 {
		distributedLocker.tickerDuration = DefaultTicker
	}
	return distributedLocker
}

func (l *DistributedLocker) WaitSignal(ctx context.Context, abortChan SyncChannel[bool], key, value string) {
	lock, err := l.redisLock.Obtain(ctx, key, l.ttl, nil)
	if err != nil {
		abortChan.Write(false)
		return
	}
	err = l.redisClient.Set(ctx, key, value, l.ttl).Err()
	if err != nil {
		abortChan.Write(false)
		return
	}
	lock.Release(context.Background())
	ticker := time.NewTicker(l.tickerDuration)
	go func() {
		for {
			select {
			case <-ctx.Done():
				abortChan.Write(false)
				return
			case <-ticker.C:
				result, err := l.redisClient.Get(ctx, key).Result()
				if err != nil || result == "" {
					abortChan.Write(false)
					return
				}
				if result != value {
					abortChan.Write(true)
					return
				}
			}
		}
	}()
}
