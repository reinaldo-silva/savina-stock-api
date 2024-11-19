package sale

import (
	"time"

	"github.com/reinaldo-silva/savina-stock/internal/domain/sale_item"
)

type Sale struct {
	ID           uint                 `gorm:"primaryKey;autoIncrement" json:"id"`
	Discount     float64              `gorm:"type:decimal(10,2);default:0" json:"discount"`
	TotalAmount  float64              `gorm:"type:decimal(10,2);not null" json:"total_amount"`
	CreatedAt    time.Time            `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time            `gorm:"autoUpdateTime" json:"updated_at"`
	SaleProducts []sale_item.SaleItem `gorm:"foreignKey:SaleID" json:"sale_products"`
}
