package product_image

import (
	"bytes"
	"fmt"

	"github.com/reinaldo-silva/savina-stock/internal/domain/image_service"
)

type ImageUseCase struct {
	imageService *image_service.ImageService
	repo         ImageRepository
}

func NewImageUseCase(ser *image_service.ImageService, repo ImageRepository) *ImageUseCase {
	return &ImageUseCase{
		imageService: ser,
		repo:         repo,
	}
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
	_, err := uc.repo.FindImageByPublicIdAndProductSlug(uuid, slug)
	if err != nil {
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
