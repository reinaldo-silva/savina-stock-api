package product

import (
	"time"

	"github.com/segmentio/ksuid"
)

type Product struct {
	ID          uint      `gorm:"primaryKey;autoIncrement"`
	Name        string    `gorm:"type:varchar(100);not null"`
	Slug        string    `gorm:"type:varchar(150);unique;not null"`
	Description string    `gorm:"type:text"`
	Price       float64   `gorm:"type:decimal(10,2);not null"`
	Cost        float64   `gorm:"type:decimal(10,2);"`
	Stock       int       `gorm:"not null"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

type ProductRepository interface {
	GetAll() ([]Product, error)
	Create(product Product) error
}

func GenerateSlug() string {
	id := ksuid.New().String()
	return id[:8]
}
