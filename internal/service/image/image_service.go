package image_service

import (
	"github.com/reinaldo-silva/savina-stock/internal/domain/image"
)

type ImageService struct {
	provider image.ImageProvider
}

func NewImageService(provider image.ImageProvider) *ImageService {
	return &ImageService{provider}
}

func (se *ImageService) GetImage(publicID string) (string, error) {
	return se.provider.GetImage(publicID)
}

func (se *ImageService) Upload(filePath string) (string, error) {
	return se.provider.UploadImage(filePath)
}
