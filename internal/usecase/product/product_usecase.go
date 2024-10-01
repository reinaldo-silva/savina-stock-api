package usecase_product

import (
	"errors"
	"fmt"
	"strings"

	"github.com/reinaldo-silva/savina-stock/internal/domain/image"
	product "github.com/reinaldo-silva/savina-stock/internal/domain/product"
)

type ProductUseCase struct {
	repo      product.ProductRepository
	imageRepo image.ImageRepository
}

func NewProductUseCase(repo product.ProductRepository, imageRepo image.ImageRepository) *ProductUseCase {
	return &ProductUseCase{repo: repo, imageRepo: imageRepo}
}

func (uc *ProductUseCase) GetAll() ([]product.Product, error) {
	return uc.repo.GetAll()
}

func (uc *ProductUseCase) Create(p product.Product) (*product.Product, error) {

	if strings.TrimSpace(p.Slug) == "" {
		p.Slug = product.GenerateSlug()
	}

	if strings.TrimSpace(p.Name) == "" {
		return nil, errors.New("product name cannot be empty")
	}

	err := uc.repo.Create(p)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (uc *ProductUseCase) GetBySlug(slug string) (*product.Product, error) {
	return uc.repo.FindBySlug(slug)
}

func (uc *ProductUseCase) Delete(slug string) error {
	err := uc.repo.DeleteBySlug(slug)
	if err != nil {
		return err
	}
	return nil
}

func (uc *ProductUseCase) Update(slug string, updatedProduct product.Product) (product.Product, error) {

	product, err := uc.repo.UpdateBySlug(slug, updatedProduct)
	if err != nil {
		return product, err
	}
	return product, nil
}

func (uc *ProductUseCase) AddImagesToProduct(slug string, imageURLs []image.UploadedImage) error {
	product, err := uc.repo.FindBySlug(slug)
	if err != nil {
		return err
	}

	if len(product.Images)+len(imageURLs) > 5 {
		return fmt.Errorf("a product can have a maximum of 5 images")
	}

	fmt.Println(product.ID, imageURLs)

	err = uc.imageRepo.CreateManyImages(product.ID, imageURLs)
	if err != nil {
		return err
	}

	return nil
}

func (uc *ProductUseCase) GetProductImages(productID uint) ([]image.ProductImage, error) {
	images, err := uc.imageRepo.FindByProductID(productID)
	if err != nil {
		return nil, err
	}
	return images, nil
}
