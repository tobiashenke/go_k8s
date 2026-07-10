package items

import "fmt"

type ItemRepository interface {
	Save(item Item) error
	Get(ID string) (*Item, error)
	GetAll() ([]Item, error)
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

func (r *InMemoryItemRepository) Get(ID string) (*Item, error) {
	item, ok := r.repo[ID]
	if !ok {
		return nil, fmt.Errorf("Item does not exist")
	}
	return &item, nil
}

func (r *InMemoryItemRepository) GetAll() ([]Item, error) {
	repoList := make([]Item, 0, len(r.repo))
	for _, v := range r.repo {
		repoList = append(repoList, v)
	}
	return repoList, nil
}
