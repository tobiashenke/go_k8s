package items

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteItemRepository struct {
	db *sql.DB
}

func NewSQLiteItemRepository() (*SQLiteItemRepository, error) {
	db, err := sql.Open("sqlite3", "items.db")
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS items (
		id   TEXT PRIMARY KEY,
		name TEXT NOT NULL
	)`)
	if err != nil {
		return nil, err
	}
	return &SQLiteItemRepository{db: db}, nil
}

func (r *SQLiteItemRepository) Save(item Item) error {
	_, err := r.db.Exec("INSERT OR REPLACE INTO items (id, name) VALUES (?, ?)", item.ID, item.Name)
	if err != nil {
		return err
	}
	return nil
}

func (r *SQLiteItemRepository) Get(id string) (*Item, error) {
	var item Item
	err := r.db.QueryRow("SELECT id, name FROM items WHERE id = ?", id).Scan(&item.ID, &item.Name)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *SQLiteItemRepository) GetAll() ([]Item, error) {
	rows, err := r.db.Query("SELECT id, name FROM items")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var itemList []Item
	for rows.Next() {
		var item Item
		err := rows.Scan(&item.ID, &item.Name)
		if err != nil {
			return nil, err
		}
		itemList = append(itemList, item)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return itemList, nil
}

func (r *SQLiteItemRepository) Delete(id string) error {
	_, err := r.db.Exec("DELETE FROM items WHERE id = ?", id)
	return err
}
