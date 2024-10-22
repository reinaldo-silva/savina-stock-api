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

func (r *GormImageRepository) CreateManyImages(productID uint, imageURLs []image.UploadedImage) error {
	for _, url := range imageURLs {
		image := image.ProductImage{
			ProductID: productID,
			ImageURL:  url.URL,
			PublicID:  url.PublicID,
			IsCover:   false,
		}
		if err := r.db.Create(&image).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *GormImageRepository) FindByProductID(productID uint) ([]image.ProductImage, error) {
	var images []image.ProductImage
	if err := r.db.Where("product_id = ?", productID).Find(&images).Error; err != nil {
		return nil, err
	}
	return images, nil
}

func (r *GormImageRepository) FindByPublicID(publicID string) (*image.ProductImage, error) {
	var img image.ProductImage
	if err := r.db.Where("public_id = ?", publicID).First(&img).Error; err != nil {
		return nil, err
	}
	return &img, nil
}
