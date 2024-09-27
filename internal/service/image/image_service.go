package service

import (
	domain "github.com/reinaldo-silva/savina-stock/internal/domain/image"
)

type ImageService struct {
	provider domain.ImageProvider
}

func NewImageService(provider domain.ImageProvider) *ImageService {
	return &ImageService{provider}
}

func (se *ImageService) GetImage(publicID string) (string, error) {
	return se.provider.GetImage(publicID)
}

func (se *ImageService) Upload(filePath string) (string, error) {
	return se.provider.UploadImage(filePath)
}
