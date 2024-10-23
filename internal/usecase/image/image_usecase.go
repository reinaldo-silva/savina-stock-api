package usecase_image

import (
	"bytes"
	"fmt"

	"github.com/reinaldo-silva/savina-stock/internal/domain/image"
	product_domain "github.com/reinaldo-silva/savina-stock/internal/domain/product"
	image_service "github.com/reinaldo-silva/savina-stock/internal/service/image"
)

type ImageUseCase struct {
	imageService *image_service.ImageService
	repo         image.ImageRepository
	productRepo  product_domain.ProductRepository
}

func NewImageUseCase(ser *image_service.ImageService, repo image.ImageRepository, productRepo product_domain.ProductRepository) *ImageUseCase {
	return &ImageUseCase{
		imageService: ser,
		repo:         repo,
		productRepo:  productRepo}
}

func (uc *ImageUseCase) GetImage(publicID string) (*bytes.Buffer, string, error) {
	img, err := uc.repo.FindByPublicID(publicID)
	if err != nil {

		return nil, "", err
	}

	if img == nil {
		return nil, "", fmt.Errorf("imagem com PublicID %s não encontrada", publicID)
	}

	return uc.imageService.Download(publicID)
}

func (uc *ImageUseCase) DeleteImage(uuid string) error {
	image, err := uc.repo.FindByPublicID(uuid)
	if err != nil {
		return fmt.Errorf("erro ao buscar imagem com UUID %s: %v", uuid, err)
	}

	if image == nil {
		return fmt.Errorf("imagem com UUID %s não encontrada", uuid)
	}

	err = uc.imageService.DeleteImage(image.PublicID)
	if err != nil {
		return fmt.Errorf("erro ao deletar imagem %s do provedor: %v", image.PublicID, err)
	}

	err = uc.repo.DeleteImage(uuid)
	if err != nil {
		return fmt.Errorf("erro ao deletar imagem %s do banco de dados: %v", uuid, err)
	}

	return nil
}

func (uc *ImageUseCase) SetImageAsCover(uuid string, slug string) error {
	image, err := uc.repo.FindByPublicID(uuid)
	if err != nil {
		return fmt.Errorf("erro ao buscar imagem com UUID %s: %v", uuid, err)
	}

	if image == nil {
		return fmt.Errorf("imagem com UUID %s não encontrada", uuid)
	}

	product, err := uc.productRepo.FindBySlug(slug)
	if err != nil || product == nil {
		return fmt.Errorf("produto com slug %s não encontrado", slug)
	}

	if image.ProductID != product.ID {
		return fmt.Errorf("imagem com UUID %s não pertence ao produto com slug %s", uuid, slug)
	}

	err = uc.repo.ResetCover(slug)
	if err != nil {
		return fmt.Errorf("erro ao resetar imagem de capa anterior: %v", err)
	}

	err = uc.repo.SetImageAsCover(uuid)
	if err != nil {
		return fmt.Errorf("erro ao definir imagem %s como capa: %v", uuid, err)
	}

	return nil
}
