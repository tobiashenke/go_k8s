package items

import (
	"context"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing"
)

type SQLiteItemRepository struct {
	db *gorm.DB
}

func NewSQLiteItemRepository() (*SQLiteItemRepository, error) {
	db, err := gorm.Open(sqlite.Open("items.db"), &gorm.Config{})
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
	return &SQLiteItemRepository{db: db}, nil
}

func (r *SQLiteItemRepository) Save(ctx context.Context, item Item) error {
	g := r.db.Save(&item)
	if g.Error != nil {
		return g.Error
	}
	return nil
}

func (r *SQLiteItemRepository) Get(ctx context.Context, id string) (*Item, error) {
	var item Item
	g := r.db.First(&item, id)
	if g.Error != nil {
		return nil, g.Error
	}
	return &item, nil
}

func (r *SQLiteItemRepository) GetAll(ctx context.Context) ([]Item, error) {
	var itemList []Item
	g := r.db.Find(&itemList)
	if g.Error != nil {
		return nil, g.Error
	}
	return itemList, nil
}

func (r *SQLiteItemRepository) Delete(ctx context.Context, id string) error {
	g := r.db.Delete(&Item{}, id)
	return g.Error
}
