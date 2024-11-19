package image_provider

import "bytes"

type Implementation interface {
	UploadImage(filePath string) (string, error)
	DownloadImage(uuid string) (*bytes.Buffer, string, error)
	DeleteImage(uuid string) error
}
