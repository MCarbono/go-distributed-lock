package locker

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type LockManager struct {
	client redis.Client
}

func NewLockManager(client redis.Client) LockManager {
	return LockManager{
		client: client,
	}
}

func (l LockManager) AcquireLock(ctx context.Context, key string) (bool, error) {
	lock, err := l.client.SetNX(ctx, key, "1", 10*time.Second).Result()
	if err != nil {
		return false, fmt.Errorf("error acquiring lock: %w", err)
	}
	return lock, nil
}

func (l LockManager) ReleaseLock(ctx context.Context, key string) error {
	_, err := l.client.Del(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("error trying to unlock: %w", err)
	}
	return nil
}

func (l LockManager) ExistsLock(ctx context.Context, key string) (bool, error) {
	exists, err := l.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("error trying to check for key %s: %w", key, err)
	}
	return exists > 0, nil
}
