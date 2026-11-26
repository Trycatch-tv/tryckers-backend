package storage

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3API define las operaciones mínimas del cliente AWS S3 para facilitar pruebas unitarias.
type S3API interface {
	PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
	GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
	DeleteObject(ctx context.Context, params *s3.DeleteObjectInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectOutput, error)
}

type S3Storage struct {
	Client        S3API
	Bucket        string
	PublicBaseURL string
}

func NewS3Storage(ctx context.Context, region, bucket, publicBaseURL string) (*S3Storage, error) {
	if strings.TrimSpace(bucket) == "" {
		return nil, fmt.Errorf("S3_BUCKET es requerido cuando STORAGE_DRIVER=s3")
	}
	if strings.TrimSpace(region) == "" {
		region = "us-west-2"
	}

	awsCfg, err := awsconfig.LoadDefaultConfig(ctx, awsconfig.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("no se pudo cargar configuración AWS SDK: %w", err)
	}

	client := s3.NewFromConfig(awsCfg)
	return &S3Storage{
		Client:        client,
		Bucket:        strings.TrimSpace(bucket),
		PublicBaseURL: strings.TrimRight(strings.TrimSpace(publicBaseURL), "/"),
	}, nil
}

func (s *S3Storage) SaveImage(ctx context.Context, fileHeader *multipart.FileHeader, options SaveOptions) (string, error) {
	img, err := ValidateAndReadImage(fileHeader, options)
	if err != nil {
		return "", err
	}

	objectKey := GenerateObjectKey(options.Kind, options.OwnerID, img.Extension)

	_, err = s.Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:       aws.String(s.Bucket),
		Key:          aws.String(objectKey),
		Body:         NewBytesReader(img.Data),
		ContentType:  aws.String(img.ContentType),
		CacheControl: aws.String("public, max-age=31536000, immutable"),
	})
	if err != nil {
		return "", fmt.Errorf("error al subir archivo a S3: %w", err)
	}

	return objectKey, nil
}

func (s *S3Storage) GetObject(ctx context.Context, objectKey string) (io.ReadCloser, string, error) {
	if !IsManagedObjectKey(objectKey) {
		return nil, "", fmt.Errorf("object_key invalido")
	}

	output, err := s.Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		return nil, "", fmt.Errorf("error al obtener objeto desde S3: %w", err)
	}

	contentType := "application/octet-stream"
	if output.ContentType != nil && *output.ContentType != "" {
		contentType = *output.ContentType
	}

	return output.Body, contentType, nil
}

func (s *S3Storage) DeleteObject(ctx context.Context, objectKey string) error {
	if !IsManagedObjectKey(objectKey) {
		return nil
	}

	_, err := s.Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		return fmt.Errorf("error al eliminar objeto en S3: %w", err)
	}
	return nil
}

func (s *S3Storage) ResolveURL(objectKey string) string {
	return ResolveMediaURL(s.PublicBaseURL, objectKey)
}
