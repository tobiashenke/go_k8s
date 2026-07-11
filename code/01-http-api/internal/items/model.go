package items

type Item struct {
	ID   string `gorm:"primaryKey"`
	Name string `gorm:"not null"`
}
