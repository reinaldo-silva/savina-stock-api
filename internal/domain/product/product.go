package product

import (
	"context"
	"time"

	"github.com/reinaldo-silva/savina-stock/internal/domain/category"
	"github.com/reinaldo-silva/savina-stock/internal/domain/product_image"

	"github.com/segmentio/ksuid"
)

type Product struct {
	ID          uint                         `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string                       `gorm:"type:varchar(100);not null" json:"name"`
	Slug        string                       `gorm:"type:varchar(150);unique;not null" json:"slug"`
	Description string                       `gorm:"type:text" json:"description"`
	Price       float64                      `gorm:"type:decimal(10,2);not null" json:"price"`
	Cost        float64                      `gorm:"type:decimal(10,2);" json:"cost"`
	Stock       int                          `gorm:"not null" json:"stock"`
	Available   bool                         `gorm:"not null;default:false" json:"available"`
	Images      []product_image.ProductImage `gorm:"foreignKey:ProductID" json:"images"`
	Categories  []category.Category          `gorm:"many2many:product_categories;" json:"categories"`
	CreatedAt   time.Time                    `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time                    `gorm:"autoUpdateTime" json:"updated_at"`
}

type ProductRepository interface {
	GetAll(
		ctx context.Context,
		page int,
		pageSize int,
		nameFilter string,
		categoryIDs []uint,
		onlyAvailable bool) ([]Product, int64, error)
	Create(product Product) error
	FindBySlug(slug string) (*Product, error)
	DeleteBySlug(productID uint) error
	UpdateBySlug(slug string, updatedProduct Product) (Product, error)
	ClearProductCategories(productID uint) error
	UpdateProductCategories(product *Product) error
	SwitchAvailable(product Product) error
}

type ProductResponse struct {
	ID          uint                         `json:"id"`
	Name        string                       `json:"name"`
	Slug        string                       `json:"slug"`
	Description string                       `json:"description"`
	Price       float64                      `json:"price"`
	Stock       int                          `json:"stock"`
	Images      []product_image.ProductImage `json:"images"`
	Categories  []category.Category          `json:"categories"`
}

func (p *Product) ToResponse() *ProductResponse {
	return &ProductResponse{
		ID:          p.ID,
		Name:        p.Name,
		Slug:        p.Slug,
		Description: p.Description,
		Price:       p.Price,
		Stock:       p.Stock,
		Images:      p.Images,
		Categories:  p.Categories,
	}
}

// func (p *Product) BeforeUpdate(tx *gorm.DB) (err error) {
// 	var oldProduct Product
// 	if err := tx.Unscoped().First(&oldProduct, p.ID).Error; err != nil {
// 		return err
// 	}

// 	userID, err := utils.GetCurrentUserID(tx)
// 	fmt.Println(userID)

// 	if err != nil {
// 		return err
// 	}

// 	oldValue := make(map[string]interface{})
// 	newValue := make(map[string]interface{})

// 	if oldProduct.Name != p.Name {
// 		oldValue["name"] = oldProduct.Name
// 		newValue["name"] = p.Name
// 	}
// 	if oldProduct.Price != p.Price {
// 		oldValue["price"] = oldProduct.Price
// 		newValue["price"] = p.Price
// 	}
// 	if oldProduct.Stock != p.Stock {
// 		oldValue["stock"] = oldProduct.Stock
// 		newValue["stock"] = p.Stock
// 	}

// 	audit := product_audit.ProductAudit{
// 		ProductID:   p.ID,
// 		UserID:      userID,
// 		Action:      "updated",
// 		OldValue:    utils.ToJSON(oldValue),
// 		NewValue:    utils.ToJSON(newValue),
// 		Description: "Product updated",
// 	}

// 	if err := tx.Create(&audit).Error; err != nil {
// 		return err
// 	}

// 	return nil
// }

func GenerateSlug() string {
	id := ksuid.New().String()
	return id[:8]
}
