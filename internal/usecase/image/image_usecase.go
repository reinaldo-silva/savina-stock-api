package usecase_image

import (
	"bytes"

	image_service "github.com/reinaldo-silva/savina-stock/internal/service/image"
)

type ImageUseCase struct {
	imageService *image_service.ImageService
}

func NewImageUseCase(ser *image_service.ImageService) *ImageUseCase {
	return &ImageUseCase{imageService: ser}
}

func (uc *ImageUseCase) GetImage(uuid string) (*bytes.Buffer, string, error) {
	return uc.imageService.Download(uuid)
}
