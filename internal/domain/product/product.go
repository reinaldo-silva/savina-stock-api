package domain

type Product struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"not null"`
	Description string
	Price       float64 `gorm:"not null"`
	Stock       int     `gorm:"not null"`
}
type ProductRepository interface {
	GetAll() ([]Product, error)
	Create(product Product) error
}
