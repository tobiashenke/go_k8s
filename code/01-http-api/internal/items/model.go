package items

type Item struct {
	ID   string `gorm:"primaryKey" validate:"required"`
	Name string `gorm:"not null" validate:"required"`
}
