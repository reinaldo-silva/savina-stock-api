package usecase

import (
	"errors"
	"fmt"
	"mime/multipart"
	"strings"

	product "github.com/reinaldo-silva/savina-stock/internal/domain/product"
	image_service "github.com/reinaldo-silva/savina-stock/internal/service/image"
)

type ProductUseCase struct {
	repo         product.ProductRepository
	imageService *image_service.ImageService
}

func NewProductUseCase(repo product.ProductRepository, cs *image_service.ImageService) *ProductUseCase {
	return &ProductUseCase{repo: repo, imageService: cs}
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

func (pu *ProductUseCase) UploadProductImage(productID string, file multipart.File) (string, error) {
	url, err := pu.imageService.Upload("")
	if err != nil {
		return "", fmt.Errorf("failed to upload image: %v", err)
	}

	return url, nil
}
