package s3_provider

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	"github.com/reinaldo-silva/savina-stock/config"
	"github.com/reinaldo-silva/savina-stock/internal/domain/image"
)

type S3Provider struct {
	Client *s3.S3
	Bucket string
}

func NewS3Provider(cfg config.S3Config) (image.ImageProvider, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(cfg.Region),
	})
	if err != nil {
		return nil, fmt.Errorf("could not create AWS session: %v", err)
	}

	return &S3Provider{
		Client: s3.New(sess),
		Bucket: cfg.BucketName,
	}, nil
}

func (sp *S3Provider) UploadImage(filePath string) (string, string, error) {
	ctx := context.Background()

	file, err := os.Open(filePath)

	if err != nil {
		return "", "", fmt.Errorf("could not open file: %v", err)
	}
	defer file.Close()

	// Gera um UUID para o nome do arquivo
	fileID := uuid.New().String()
	fileExt := filepath.Ext(filePath)
	s3Key := fmt.Sprintf("%s%s", fileID, fileExt)

	// Define as opções de upload
	_, err = sp.Client.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(sp.Bucket),
		Key:         aws.String(s3Key),
		Body:        file,
		ContentType: aws.String("image/jpeg"), // Ajuste o tipo de conteúdo conforme necessário
	})
	fmt.Println(err)
	if err != nil {
		return "", "", fmt.Errorf("could not upload image to S3: %v", err)
	}

	// Gera a URL do objeto
	url := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", sp.Bucket, s3Key)
	return url, fileID, nil
}

func (sp *S3Provider) GetImage(fileID string) (string, error) {
	// Gera a URL do objeto com o UUID (sem extensão)
	s3Key := fmt.Sprintf("%s.jpg", fileID) // Ajuste a extensão conforme necessário
	url := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", sp.Bucket, s3Key)
	return url, nil
}
