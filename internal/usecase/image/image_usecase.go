package usecase_image

import (
	"bytes"
	"fmt"

	"github.com/reinaldo-silva/savina-stock/internal/domain/image"
	image_service "github.com/reinaldo-silva/savina-stock/internal/service/image"
)

type ImageUseCase struct {
	imageService *image_service.ImageService
	repo         image.ImageRepository
}

func NewImageUseCase(ser *image_service.ImageService, repo image.ImageRepository) *ImageUseCase {
	return &ImageUseCase{imageService: ser, repo: repo}
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
