package cloudinary_provider

import (
	"bytes"
	"context"
	"fmt"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/reinaldo-silva/savina-stock/config"
	"github.com/reinaldo-silva/savina-stock/internal/domain/image"
)

type CloudinaryProvider struct {
	Client *cloudinary.Cloudinary
}

func NewCloudinaryProvider(cfg config.CloudinaryConfig) (image.ImageProvider, error) {
	cld, err := cloudinary.NewFromParams(cfg.CloudName, cfg.APIKey, cfg.APISecret)
	if err != nil {
		return nil, fmt.Errorf("could not initialize Cloudinary: %v", err)
	}

	return &CloudinaryProvider{Client: cld}, nil
}

func (cp *CloudinaryProvider) UploadImage(filePath string) (string, error) {
	ctx := context.Background()

	resp, err := cp.Client.Upload.Upload(ctx, filePath, uploader.UploadParams{})
	if err != nil {
		return "", fmt.Errorf("could not upload image: %v", err)
	}

	return resp.PublicID, nil
}

func (cp *CloudinaryProvider) GetImage(publicID string) (string, error) {

	img, err := cp.Client.Image(publicID)
	if err != nil {
		return "", fmt.Errorf("could not get image: %v", err)
	}

	url, err := img.String()
	if err != nil {
		return "", fmt.Errorf("could not generate image URL: %v", err)
	}

	return url, nil
}

func (sp *CloudinaryProvider) DownloadImage(uuid string) (*bytes.Buffer, string, error) {
	return nil, "", nil
}

func (sp *CloudinaryProvider) DeleteImage(uuid string) error {
	return nil
}
