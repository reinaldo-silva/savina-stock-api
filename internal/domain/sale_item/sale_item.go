package sale_item

import (
	"time"

	"github.com/reinaldo-silva/savina-stock/internal/domain/product"
)

type SaleItem struct {
	ID        uint            `gorm:"primaryKey;autoIncrement" json:"id"`
	SaleID    uint            `gorm:"not null;index" json:"sale_id"`
	ProductID uint            `gorm:"not null;index" json:"product_id"`
	Quantity  int             `gorm:"not null" json:"quantity"`
	UnitPrice float64         `gorm:"type:decimal(10,2);not null" json:"unit_price"`
	SubTotal  float64         `gorm:"type:decimal(10,2);not null" json:"subtotal"`
	CreatedAt time.Time       `gorm:"autoCreateTime" json:"created_at"`
	Product   product.Product `gorm:"foreignKey:ProductID" json:"product"`
}
