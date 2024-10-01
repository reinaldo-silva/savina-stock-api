package image_repository

import (
	"github.com/reinaldo-silva/savina-stock/internal/domain/image"
	"gorm.io/gorm"
)

type GormImageRepository struct {
	db *gorm.DB
}

func NewGormImageRepository(db *gorm.DB) *GormImageRepository {
	return &GormImageRepository{db: db}
}

func (r *GormImageRepository) CreateManyImages(productID uint, imageURLs []string) error {
	for _, url := range imageURLs {
		image := image.ProductImage{
			ProductID: productID,
			ImageURL:  url,
			IsCover:   false,
		}
		if err := r.db.Create(&image).Error; err != nil {
			return err
		}
	}
	return nil
}
