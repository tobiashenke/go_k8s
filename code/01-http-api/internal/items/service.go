package items

type ItemService struct {
	itemRepository ItemRepository
}

func NewItemService(itemRepo ItemRepository) *ItemService {
	return &ItemService{
		itemRepository: itemRepo,
	}
}

func (i *ItemService) Get(Id string) (*Item, error) {
	return i.itemRepository.Get(Id)
}

func (i *ItemService) Save(item Item) error {
	return i.itemRepository.Save(item)
}

func (i *ItemService) GetAll() ([]Item, error) {
	return i.itemRepository.GetAll()
}
