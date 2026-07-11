package items

import "fmt"

type ItemRepository interface {
	Save(item Item) error
	Get(id string) (*Item, error)
	GetAll() ([]Item, error)
	Delete(id string) error
}

type InMemoryItemRepository struct {
	repo map[string]Item
}

func NewInMemoryItemRepository() *InMemoryItemRepository {
	return &InMemoryItemRepository{repo: make(map[string]Item)}
}

func (r *InMemoryItemRepository) Save(item Item) error {
	r.repo[item.ID] = item
	return nil
}

func (r *InMemoryItemRepository) Get(id string) (*Item, error) {
	item, ok := r.repo[id]
	if !ok {
		return nil, fmt.Errorf("item does not exist")
	}
	return &item, nil
}

func (r *InMemoryItemRepository) GetAll() ([]Item, error) {
	itemList := make([]Item, 0, len(r.repo))
	for _, v := range r.repo {
		itemList = append(itemList, v)
	}
	return itemList, nil
}

func (r *InMemoryItemRepository) Delete(id string) error {
	delete(r.repo, id)
	return nil
}
