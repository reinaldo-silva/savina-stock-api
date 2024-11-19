package product_audit

import (
	"time"
)

type ProductAudit struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ProductID   uint      `gorm:"not null;index" json:"product_id"`
	UserID      uint      `gorm:"not null;index" json:"user_id"`
	Action      string    `gorm:"type:varchar(50);not null" json:"action"` // Ex: "created", "updated"
	OldValue    string    `gorm:"type:text" json:"old_value"`
	NewValue    string    `gorm:"type:text" json:"new_value"`
	Description string    `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
}
