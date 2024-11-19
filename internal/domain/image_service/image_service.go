package image_service

import (
	"bytes"

	"github.com/reinaldo-silva/savina-stock/internal/infrastructure/image_provider"
)

type ImageService struct {
	provider image_provider.Implementation
}

func NewImageService(provider image_provider.Implementation) *ImageService {
	return &ImageService{provider}
}

func (se *ImageService) Upload(filePath string) (string, error) {
	return se.provider.UploadImage(filePath)
}

func (se *ImageService) Download(uuid string) (*bytes.Buffer, string, error) {
	return se.provider.DownloadImage(uuid)
}

func (se *ImageService) DeleteImage(uuid string) error {
	return se.provider.DeleteImage(uuid)
}
