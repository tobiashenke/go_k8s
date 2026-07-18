package items

import (
	"context"
	"fmt"
)

type ItemRepository interface {
	Save(ctx context.Context, item Item) error
	Get(ctx context.Context, id string) (*Item, error)
	GetAll(ctx context.Context) ([]Item, error)
	Delete(ctx context.Context, id string) error
}

type InMemoryItemRepository struct {
	repo map[string]Item
}

func NewInMemoryItemRepository() *InMemoryItemRepository {
	return &InMemoryItemRepository{repo: make(map[string]Item)}
}

func (r *InMemoryItemRepository) Save(ctx context.Context, item Item) error {
	r.repo[item.ID] = item
	return nil
}

func (r *InMemoryItemRepository) Get(ctx context.Context, id string) (*Item, error) {
	item, ok := r.repo[id]
	if !ok {
		return nil, fmt.Errorf("item does not exist")
	}
	return &item, nil
}

func (r *InMemoryItemRepository) GetAll(ctx context.Context) ([]Item, error) {
	itemList := make([]Item, 0, len(r.repo))
	for _, v := range r.repo {
		itemList = append(itemList, v)
	}
	return itemList, nil
}

func (r *InMemoryItemRepository) Delete(ctx context.Context, id string) error {
	delete(r.repo, id)
	return nil
}
