package image

import (
	"time"
)

type ProductImage struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ProductID uint      `gorm:"not null" json:"product_id"`
	ImageURL  string    `gorm:"type:varchar(255);not null" json:"image_url"`
	PublicID  string    `gorm:"type:varchar(255);not null" json:"public_id"`
	IsCover   bool      `gorm:"default:false" json:"is_cover"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (ProductImage) TableName() string {
	return "product_images"
}

type ImageRepository interface {
	CreateManyImages(productID uint, imageURLs []UploadedImage) error
}

type ImageProvider interface {
	UploadImage(filePath string) (string, string, error)
	GetImage(publicID string) (string, error)
}

type UploadedImage struct {
	URL      string `json:"url"`
	PublicID string `json:"public_id"`
}
