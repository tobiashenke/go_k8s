package items

import (
	"context"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing"
)

type PostgresItemRepository struct {
	db *gorm.DB
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
	return &PostgresItemRepository{db: db}, nil
}

func (p *PostgresItemRepository) Save(ctx context.Context, item Item) error {
	g := p.db.Save(&item)
	if g.Error != nil {
		return g.Error
	}
	return nil
}

func (p *PostgresItemRepository) Get(ctx context.Context, id string) (*Item, error) {
	var item Item
	g := p.db.First(&item, id)
	if g.Error != nil {
		return nil, g.Error
	}
	return &item, nil
}

func (p *PostgresItemRepository) GetAll(ctx context.Context) ([]Item, error) {
	var itemList []Item
	g := p.db.Find(&itemList)
	if g.Error != nil {
		return nil, g.Error
	}
	return itemList, nil
}

func (p *PostgresItemRepository) Delete(ctx context.Context, id string) error {
	g := p.db.Delete(&Item{}, id)
	return g.Error
}
