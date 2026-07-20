package items

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
	"github.com/sony/gobreaker/v2"
)

type ItemCache struct {
	client *redis.Client
	cb     *gobreaker.CircuitBreaker[any]
}

func NewItemCache(RedisAddress string) (*ItemCache, error) {
	client := &ItemCache{
		client: redis.NewClient(&redis.Options{Addr: RedisAddress}),
		cb: gobreaker.NewCircuitBreaker[any](gobreaker.Settings{
			Name:        "redis",
			MaxRequests: 1,
			Interval:    60 * time.Second,
			Timeout:     30 * time.Second,
			ReadyToTrip: func(counts gobreaker.Counts) bool {
				return counts.ConsecutiveFailures > 5
			},
		}),
	}
	err := redisotel.InstrumentTracing(client.client)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (c *ItemCache) Set(ctx context.Context, item Item) error {
	bytes, err := json.Marshal(item)
	if err != nil {
		return err
	}
	_, err = c.cb.Execute(func() (any, error) {
		return nil, c.client.Set(ctx, item.ID, bytes, time.Minute*5).Err()
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *ItemCache) Get(ctx context.Context, id string) (*Item, error) {
	result, err := c.cb.Execute(func() (any, error) {
		return c.client.Get(ctx, id).Bytes()
	})
	if err != nil {
		return nil, err
	}
	var item Item
	err = json.Unmarshal(result.([]byte), &item)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (c *ItemCache) Delete(ctx context.Context, id string) error {
	_, err := c.cb.Execute(func() (any, error) {
		return nil, c.client.Del(ctx, id).Err()
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *ItemCache) SetResponse(ctx context.Context, key string, data []byte, ttl time.Duration) error {
	r := c.client.Set(ctx, key, data, ttl)
	if r.Err() != nil {
		return r.Err()
	}
	return nil
}

func (c *ItemCache) GetResponse(ctx context.Context, key string) ([]byte, error) {
	return c.client.Get(ctx, key).Bytes()
}

func (c *ItemCache) ExceededRateLimit(ctx context.Context, ip string, limit int, window time.Duration) (bool, error) {
	windowStart := time.Now().Add(-window).UnixNano()
	z := c.client.ZRemRangeByScore(ctx, ip, "0", fmt.Sprintf("%d", windowStart))
	if z.Err() != nil {
		return false, nil
	}
	count := c.client.ZCard(ctx, ip)
	// fail open: don't block on Redis error -> allow request to go through
	if count.Err() != nil {
		return false, nil
	}
	if count.Val() < int64(limit) {
		z = c.client.ZAdd(ctx, ip, redis.Z{
			Score:  float64(time.Now().UnixNano()),
			Member: time.Now().UnixNano(),
		})
		if z.Err() != nil {
			return false, nil
		}
		b := c.client.Expire(ctx, ip, window)
		if b.Err() != nil {
			return false, nil
		}
	}
	return count.Val() >= int64(limit), nil
}

var LuaRateLimitScript = `
redis.call('ZREMRANGEBYSCORE', KEYS[1], 0, ARGV[2])
local count = redis.call('ZCARD', KEYS[1])
if count < tonumber(ARGV[1]) then
	redis.call('ZADD', KEYS[1], ARGV[3], ARGV[3])
	redis.call('EXPIRE', KEYS[1], tonumber(ARGV[4]))
	return 0
else
	return 1
end
`

func (c *ItemCache) ExceededRateLimitWithLua(ctx context.Context, ip string, limit int, window time.Duration) (bool, error) {
	windowStart := time.Now().Add(-window).UnixNano()
	now := time.Now().UnixNano()
	cmd := c.client.Eval(ctx, LuaRateLimitScript, []string{ip}, limit, windowStart, now, int(window.Seconds()))
	return cmd.Bool()
}
