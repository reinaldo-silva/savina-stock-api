package usecase

import domain "github.com/reinaldo-silva/savina-stock/internal/domain/product"

type ProductUseCase struct {
	repo domain.ProductRepository
}

func NewProductUseCase(repo domain.ProductRepository) *ProductUseCase {
	return &ProductUseCase{repo}
}

func (uc *ProductUseCase) GetAllProducts() ([]domain.Product, error) {
	return uc.repo.GetAll()
}

func (uc *ProductUseCase) AddProduct(product domain.Product) error {
	return uc.repo.Create(product)
}
