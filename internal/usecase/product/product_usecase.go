package usecase_product

import (
	"errors"
	"fmt"
	"strings"

	category "github.com/reinaldo-silva/savina-stock/internal/domain/category"
	"github.com/reinaldo-silva/savina-stock/internal/domain/image"
	product "github.com/reinaldo-silva/savina-stock/internal/domain/product"
	image_service "github.com/reinaldo-silva/savina-stock/internal/service/image"
	"github.com/reinaldo-silva/savina-stock/utils"
)

type ProductUseCase struct {
	repo         product.ProductRepository
	categoryRepo category.CategoryRepository
	imageRepo    image.ImageRepository
	imageService *image_service.ImageService
}

func NewProductUseCase(repo product.ProductRepository, categoryRepo category.CategoryRepository, imageRepo image.ImageRepository, imageService *image_service.ImageService) *ProductUseCase {
	return &ProductUseCase{
		repo:         repo,
		categoryRepo: categoryRepo,
		imageRepo:    imageRepo,
		imageService: imageService}
}

func (uc *ProductUseCase) GetAll(page int,
	pageSize int,
	nameFilter string,
	categoryIDs []uint,
	host string) ([]product.ProductResponse, int64, error) {
	products, total, err := uc.repo.GetAll(page, pageSize, nameFilter, categoryIDs, true)
	if err != nil {
		return nil, 0, err
	}

	for i := range products {
		for j := range products[i].Images {
			products[i].Images[j].ImageURL = utils.GenerateImageURL(host, products[i].Images[j].PublicID)
		}
	}

	var productResponses []product.ProductResponse
	for _, p := range products {
		productResponses = append(productResponses, *p.ToResponse())
	}

	return productResponses, total, nil
}

func (uc *ProductUseCase) GetAllToAdmin(page int,
	pageSize int,
	nameFilter string,
	categoryIDs []uint,
	host string) ([]product.Product, int64, error) {
	products, total, err := uc.repo.GetAll(page, pageSize, nameFilter, categoryIDs, false)
	if err != nil {
		return nil, 0, err
	}

	for i := range products {
		for j := range products[i].Images {
			products[i].Images[j].ImageURL = utils.GenerateImageURL(host, products[i].Images[j].PublicID)
		}
	}

	return products, total, nil
}

func (uc *ProductUseCase) Create(p product.Product) (*product.Product, error) {

	if strings.TrimSpace(p.Slug) == "" {
		p.Slug = product.GenerateSlug()
	}

	if strings.TrimSpace(p.Name) == "" {
		return nil, errors.New("product name cannot be empty")
	}

	var categories []category.Category
	for _, category := range p.Categories {
		foundCategory, err := uc.categoryRepo.GetByID(category.ID)
		if err != nil {
			return nil, fmt.Errorf("category with ID %d does not exist", category.ID)
		}
		categories = append(categories, *foundCategory)
	}

	p.Categories = categories

	err := uc.repo.Create(p)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (uc *ProductUseCase) GetBySlug(slug string) (*product.ProductResponse, error) {
	product, err := uc.repo.FindBySlug(slug)
	if err != nil {
		return nil, err
	}

	return product.ToResponse(), nil
}

func (uc *ProductUseCase) GetBySlugToAdmin(slug string) (*product.Product, error) {
	product, err := uc.repo.FindBySlug(slug)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (uc *ProductUseCase) Delete(slug string) error {

	product, err := uc.repo.FindBySlug(slug)
	if err != nil {
		return err
	}

	if product == nil {
		return fmt.Errorf("produto com slug %s não encontrado", slug)
	}

	images, err := uc.imageRepo.FindByProductID(product.ID)
	if err != nil {
		return fmt.Errorf("erro ao buscar imagens do produto: %v", err)
	}

	for _, img := range images {
		err := uc.imageService.DeleteImage(img.PublicID)
		if err != nil {
			return fmt.Errorf("erro ao deletar imagem %s do S3: %v", img.PublicID, err)
		}
	}

	err = uc.imageRepo.DeleteByProductID(product.ID)
	if err != nil {
		return err
	}

	err = uc.repo.DeleteBySlug(product.ID)
	if err != nil {
		return err
	}

	return nil
}

func (uc *ProductUseCase) Update(slug string, updatedProduct product.Product) (*product.Product, error) {
	product, err := uc.repo.UpdateBySlug(slug, updatedProduct)
	if err != nil {
		return &product, err
	}

	return &product, nil
}

func (uc *ProductUseCase) AddImagesToProduct(slug string, imageURLs []image.UploadedImage) error {
	product, err := uc.repo.FindBySlug(slug)
	if err != nil {
		return err
	}

	if len(product.Images)+len(imageURLs) > 5 {
		return fmt.Errorf("a product can have a maximum of 5 images")
	}

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

func (uc *ProductUseCase) UpdateProductCategories(slug string, categoryIDs []int) error {
	product, err := uc.repo.FindBySlug(slug)
	if err != nil {
		return fmt.Errorf("product with slug %s not found", slug)
	}

	if product == nil {
		return fmt.Errorf("product with slug %s not found", slug)
	}

	var categories []category.Category
	for _, categoryID := range categoryIDs {
		foundCategory, err := uc.categoryRepo.GetByID(uint(categoryID))
		if err != nil {
			return fmt.Errorf("category with ID %d does not exist", categoryID)
		}
		categories = append(categories, *foundCategory)
	}

	err = uc.repo.ClearProductCategories(product.ID)
	if err != nil {
		return fmt.Errorf("failed to clear product categories: %v", err)
	}

	product.Categories = categories

	err = uc.repo.UpdateProductCategories(product)
	if err != nil {
		return fmt.Errorf("failed to update product categories: %v", err)
	}

	return nil
}
