package domain

type ImageProvider interface {
	UploadImage(filePath string) (string, error)
	GetImage(publicID string) (string, error)
}
