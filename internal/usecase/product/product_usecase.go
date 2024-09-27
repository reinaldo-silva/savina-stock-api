package usecase

import (
	"fmt"
	"mime/multipart"

	domain "github.com/reinaldo-silva/savina-stock/internal/domain/product"
	service "github.com/reinaldo-silva/savina-stock/internal/service/image"
)

type ProductUseCase struct {
	repo         domain.ProductRepository
	imageService *service.ImageService
}

func NewProductUseCase(repo domain.ProductRepository, cs *service.ImageService) *ProductUseCase {
	return &ProductUseCase{repo: repo, imageService: cs}
}

func (uc *ProductUseCase) GetAllProducts() ([]domain.Product, error) {
	return uc.repo.GetAll()
}

func (uc *ProductUseCase) AddProduct(product domain.Product) error {
	return uc.repo.Create(product)
}

func (pu *ProductUseCase) UploadProductImage(productID string, file multipart.File) (string, error) {
	url, err := pu.imageService.Upload("")
	if err != nil {
		return "", fmt.Errorf("failed to upload image: %v", err)
	}

	return url, nil
}
