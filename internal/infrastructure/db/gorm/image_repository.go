package gorm

import (
	"fmt"

	"github.com/reinaldo-silva/savina-stock/internal/domain/product"
	"github.com/reinaldo-silva/savina-stock/internal/domain/product_image"
	"gorm.io/gorm"
)

type GormImageRepository struct {
	db *gorm.DB
}

func NewGormImageRepository(db *gorm.DB) product_image.ImageRepository {
	return &GormImageRepository{db: db}
}

func (r *GormImageRepository) CreateManyImages(productID uint, imageURLs []product_image.UploadedImage) error {
	for _, url := range imageURLs {
		image := product_image.ProductImage{
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

func (r *GormImageRepository) FindByProductID(productID uint) ([]product_image.ProductImage, error) {
	var images []product_image.ProductImage
	if err := r.db.Where("product_id = ?", productID).Find(&images).Error; err != nil {
		return nil, err
	}
	return images, nil
}

func (r *GormImageRepository) FindImageByPublicIdAndProductSlug(publicID string, slug string) (*product_image.ProductImage, error) {
	var image product_image.ProductImage

	if err := r.db.Joins("JOIN products ON products.id = product_images.product_id").
		Where("product_images.public_id = ? AND products.slug = ?", publicID, slug).
		First(&image).Error; err != nil {
		return nil, err
	}

	return &image, nil
}

func (r *GormImageRepository) FindByPublicID(publicID string) (*product_image.ProductImage, error) {
	var img product_image.ProductImage
	if err := r.db.Where("public_id = ?", publicID).First(&img).Error; err != nil {
		return nil, err
	}
	return &img, nil
}

func (r *GormImageRepository) DeleteByProductID(productID uint) error {
	if err := r.db.Where("product_id = ?", productID).Delete(&product_image.ProductImage{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *GormImageRepository) DeleteImage(uuid string) error {
	return r.db.Where("public_id = ?", uuid).Delete(&product_image.ProductImage{}).Error
}

func (r *GormImageRepository) ResetCover(slug string) error {

	var productID uint
	if err := r.db.Model(&product.Product{}).Where("slug = ?", slug).Select("id").Scan(&productID).Error; err != nil {
		return fmt.Errorf("erro ao buscar produto pelo slug %s: %v", slug, err)
	}

	if productID == 0 {
		return fmt.Errorf("produto com slug %s não encontrado", slug)
	}

	return r.db.Model(&product_image.ProductImage{}).
		Where("product_id = ?", productID).
		Update("is_cover", false).Error
}

func (r *GormImageRepository) SetImageAsCover(uuid string) error {
	return r.db.Model(&product_image.ProductImage{}).Where("public_id = ?", uuid).Update("is_cover", true).Error
}
