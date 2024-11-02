package product

import (
	"time"

	category_model "github.com/reinaldo-silva/savina-stock/internal/domain/category"
	"github.com/reinaldo-silva/savina-stock/internal/domain/image"
	"github.com/segmentio/ksuid"
)

type Product struct {
	ID          uint                      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string                    `gorm:"type:varchar(100);not null" json:"name"`
	Slug        string                    `gorm:"type:varchar(150);unique;not null" json:"slug"`
	Description string                    `gorm:"type:text" json:"description"`
	Price       float64                   `gorm:"type:decimal(10,2);not null" json:"price"`
	Cost        float64                   `gorm:"type:decimal(10,2);" json:"cost"`
	Stock       int                       `gorm:"not null" json:"stock"`
	Images      []image.ProductImage      `gorm:"foreignKey:ProductID" json:"images"`
	Categories  []category_model.Category `gorm:"many2many:product_categories;" json:"categories"`
	CreatedAt   time.Time                 `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time                 `gorm:"autoUpdateTime" json:"updated_at"`
}

type ProductRepository interface {
	GetAll(
		page int,
		pageSize int,
		nameFilter string,
		categoryIDs []uint) ([]Product, int64, error)
	Create(product Product) error
	FindBySlug(slug string) (*Product, error)
	DeleteBySlug(productID uint) error
	UpdateBySlug(slug string, updatedProduct Product) (Product, error)
	ClearProductCategories(productID uint) error
	UpdateProductCategories(product *Product) error
}

func GenerateSlug() string {
	id := ksuid.New().String()
	return id[:8]
}
