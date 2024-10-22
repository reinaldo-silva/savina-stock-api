package s3_provider

import (
	"bytes"
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

func (sp *S3Provider) UploadImage(filePath string) (string, error) {
	ctx := context.Background()

	file, err := os.Open(filePath)

	if err != nil {
		return "", fmt.Errorf("could not open file: %v", err)
	}
	defer file.Close()

	fileID := uuid.New().String()
	fileExt := filepath.Ext(filePath)
	s3Key := fmt.Sprintf("%s%s", fileID, fileExt)

	_, err = sp.Client.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(sp.Bucket),
		Key:         aws.String(s3Key),
		Body:        file,
		ContentType: aws.String("image/jpeg"),
	})

	if err != nil {
		return "", fmt.Errorf("could not upload image to S3: %v", err)
	}

	return s3Key, nil
}

func (sp *S3Provider) DownloadImage(uuid string) (*bytes.Buffer, string, error) {
	ctx := context.Background()

	s3Key := uuid

	result, err := sp.Client.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: aws.String(sp.Bucket),
		Key:    aws.String(s3Key),
	})

	if err != nil {
		return nil, "", fmt.Errorf("could not download image from S3: %v", err)
	}
	defer result.Body.Close()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(result.Body)
	if err != nil {
		return nil, "", fmt.Errorf("could not read image data: %v", err)
	}

	return buf, *result.ContentType, nil
}

func (sp *S3Provider) DeleteImage(uuid string) error {
	ctx := context.Background()

	_, err := sp.Client.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(sp.Bucket),
		Key:    aws.String(uuid),
	})

	if err != nil {
		return fmt.Errorf("could not delete image from S3: %v", err)
	}

	err = sp.Client.WaitUntilObjectNotExistsWithContext(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(sp.Bucket),
		Key:    aws.String(uuid),
	})

	if err != nil {
		return fmt.Errorf("failed to confirm deletion of image: %v", err)
	}

	return nil
}
