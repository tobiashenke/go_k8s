package items

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type ItemService struct {
	r ItemRepository
	c *ItemCache
}

func NewItemService(ir ItemRepository, ic *ItemCache) *ItemService {
	return &ItemService{
		r: ir,
		c: ic,
	}
}

func (i *ItemService) Get(id string) (*Item, error) {
	item, err := i.c.Get(context.Background(), id)
	if err == nil {
		return item, err
	}
	if err == redis.Nil {
		item, err = i.r.Get(id)
		if err != nil {
			return nil, err
		}
		i.c.Set(context.Background(), *item)
		return item, err
	}
	return nil, err
}

func (i *ItemService) Save(item Item) error {
	return i.r.Save(item)
}

func (i *ItemService) GetAll() ([]Item, error) {
	return i.r.GetAll()
}
