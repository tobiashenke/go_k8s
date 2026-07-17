package items

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type ItemCache struct {
	client *redis.Client
}

func NewItemCache(RedisAddress string) *ItemCache {
	return &ItemCache{
		client: redis.NewClient(&redis.Options{Addr: RedisAddress}),
	}
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
