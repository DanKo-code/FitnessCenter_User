package localstack_usecase

import (
	"User/internal/models"
	log_c "User/pkg/logger"
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type LocalstackUseCase struct {
	client *s3.Client
	config *models.CloudConfig
}

func NewLocalstackUseCase(client *s3.Client, config *models.CloudConfig) *LocalstackUseCase {
	return &LocalstackUseCase{
		client: client,
		config: config,
	}
}

func (luc *LocalstackUseCase) PutObject(ctx context.Context, object []byte, name string) (string, error) {
	_, err := luc.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(luc.config.Bucket),
		Key:    aws.String(name),
		Body:   bytes.NewReader(object),
	})
	if err != nil {
		log_c.ErrorLogger.Printf("Failed to put object: %v", err)
		return "", err
	}

	fileURL := fmt.Sprintf("%s/%s/%s", luc.config.EndPoint, luc.config.Bucket, name)

	return fileURL, nil
}
