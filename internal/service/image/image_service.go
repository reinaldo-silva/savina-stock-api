package image_service

import (
	"bytes"

	"github.com/reinaldo-silva/savina-stock/internal/domain/image"
)

type ImageService struct {
	provider image.ImageProvider
}

func NewImageService(provider image.ImageProvider) *ImageService {
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
