package items

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

type ItemCache struct {
	client *redis.Client
}

func NewItemCache(RedisAddress string) (*ItemCache, error) {
	client := &ItemCache{
		client: redis.NewClient(&redis.Options{Addr: RedisAddress}),
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
	r := c.client.Set(ctx, item.ID, bytes, time.Minute*5)
	if r.Err() != nil {
		return r.Err()
	}
	return nil
}

func (c *ItemCache) Get(ctx context.Context, id string) (*Item, error) {
	bytes, err := c.client.Get(ctx, id).Bytes()
	if err != nil {
		return nil, err
	}
	var item Item
	err = json.Unmarshal(bytes, &item)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (c *ItemCache) Delete(ctx context.Context, id string) error {
	r := c.client.Del(ctx, id)
	if r.Err() != nil {
		return r.Err()
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
	return count.Val() > int64(limit), nil
}
