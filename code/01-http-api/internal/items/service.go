package items

import (
	"context"

	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel"
)

type ItemService struct {
	r ItemRepository
	c *ItemCache
	p *ItemPublisher
}

func NewItemService(ir ItemRepository, ic *ItemCache, ip *ItemPublisher) *ItemService {
	return &ItemService{
		r: ir,
		c: ic,
		p: ip,
	}
}

func (i *ItemService) Get(ctx context.Context, id string) (*Item, error) {
	ctx, span := otel.Tracer("http-api").Start(ctx, "service.Get")
	defer span.End()

	item, err := i.c.Get(ctx, id)
	if err == nil {
		return item, err
	}
	if err == redis.Nil {
		item, err = i.r.Get(ctx, id)
		if err != nil {
			return nil, err
		}
		i.c.Set(ctx, *item)
		return item, err
	}
	return nil, err
}

func (i *ItemService) Save(ctx context.Context, item Item) error {
	ctx, span := otel.Tracer("http-api").Start(ctx, "service.Save")
	defer span.End()

	err := i.r.Save(ctx, item)
	if err != nil {
		return err
	}
	err = i.p.PublishCreate(ctx, item)
	return err
}

func (i *ItemService) GetAll(ctx context.Context) ([]Item, error) {
	ctx, span := otel.Tracer("http-api").Start(ctx, "service.GetAll")
	defer span.End()

	return i.r.GetAll(ctx)
}

func (i *ItemService) Delete(ctx context.Context, id string) error {
	ctx, span := otel.Tracer("http-api").Start(ctx, "service.Delete")
	defer span.End()

	err := i.r.Delete(ctx, id)
	if err != nil {
		return err
	}
	err = i.c.Delete(ctx, id)
	if err != nil {
		return err
	}
	err = i.p.PublishDelete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
