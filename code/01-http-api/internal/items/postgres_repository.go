package items

import (
	"context"
	"time"

	"github.com/sony/gobreaker/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing"
)

type PostgresItemRepository struct {
	db *gorm.DB
	cb *gobreaker.CircuitBreaker[any]
}

func NewPostgresItemRepository(dsn string) (*PostgresItemRepository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&Item{})
	if err != nil {
		return nil, err
	}
	err = db.Use(tracing.NewPlugin())
	if err != nil {
		return nil, err
	}
	cb := gobreaker.NewCircuitBreaker[any](gobreaker.Settings{
		Name:        "postgres",
		MaxRequests: 1,
		Interval:    60 * time.Second,
		Timeout:     30 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures > 5
		},
	})
	return &PostgresItemRepository{db: db, cb: cb}, nil
}

func (p *PostgresItemRepository) Save(ctx context.Context, item Item) error {
	_, err := p.cb.Execute(func() (any, error) {
		return nil, p.db.Save(&item).Error
	})
	return err
}

func (p *PostgresItemRepository) Get(ctx context.Context, id string) (*Item, error) {
	result, err := p.cb.Execute(func() (any, error) {
		var item Item
		g := p.db.First(&item, id)
		return &item, g.Error
	})
	if err != nil {
		return nil, err
	}
	return result.(*Item), nil
}

func (p *PostgresItemRepository) GetAll(ctx context.Context) ([]Item, error) {
	result, err := p.cb.Execute(func() (any, error) {
		var itemList []Item
		g := p.db.Find(&itemList)
		return itemList, g.Error
	})
	if err != nil {
		return nil, err
	}
	return result.([]Item), nil
}

func (p *PostgresItemRepository) Delete(ctx context.Context, id string) error {
	_, err := p.cb.Execute(func() (any, error) {
		return nil, p.db.Delete(&Item{}, id).Error
	})
	return err
}
