package storage

import (
	"context"
	"mime"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Storage interface {
	PutPublicFile(objFile multipart.File, objKey string) (*s3.PutObjectOutput, error)
	DeletePublicFile(objKey string) (*s3.DeleteObjectOutput, error)
	GetPresignURL(objKey string) (string, error)
}

type AWSConfig struct {
	AccessKey        string `mapstructure:"AWS_ACCESS_KEY" validate:"required"`
	SecretKey        string `mapstructure:"AWS_SECRET_KEY" validate:"required"`
	Region           string `mapstructure:"AWS_REGION" validate:"required"`
	Endpoint         string `mapstructure:"AWS_S3_ENDPOINT" validate:"required"`
	BucketName       string `mapstructure:"AWS_S3_BUCKET_NAME" validate:"required"`
	PublicBucketName string `mapstructure:"AWS_S3_PUBLIC_BUCKET_NAME" validate:"required"`
}

type S3Client struct {
	s3  *s3.Client
	cfg *AWSConfig
}

func NewS3Client(AWScfg *AWSConfig) Storage {

	// configuration aws s3 Endpoint resolver
	resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...any) (aws.Endpoint, error) {
		return aws.Endpoint{
			PartitionID:       "aws",
			SigningRegion:     region,
			URL:               AWScfg.Endpoint,
			HostnameImmutable: true,
		}, nil
	})

	// configuration aws s3 main
	cfg := aws.Config{
		Region:                      AWScfg.Region,
		Credentials:                 credentials.NewStaticCredentialsProvider(AWScfg.AccessKey, AWScfg.SecretKey, ""),
		EndpointResolverWithOptions: resolver,
	}

	s3 := s3.NewFromConfig(cfg, func(o *s3.Options) { o.UsePathStyle = true })

	return &S3Client{
		s3:  s3,
		cfg: AWScfg,
	}
}

func (s *S3Client) PutPublicFile(objFile multipart.File, objKey string) (*s3.PutObjectOutput, error) {
	contentType := mime.TypeByExtension(filepath.Ext(objKey))
	output, err := s.s3.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket:      aws.String(s.cfg.PublicBucketName),
		Key:         aws.String(objKey),
		Body:        objFile,
		ContentType: aws.String(contentType),
		ACL:         "public-read",
	})

	if err != nil {
		return nil, err
	}

	return output, nil
}

func (s *S3Client) DeletePublicFile(objKey string) (*s3.DeleteObjectOutput, error) {
	result, err := s.s3.DeleteObject(context.Background(), &s3.DeleteObjectInput{
		Bucket: aws.String(s.cfg.PublicBucketName),
		Key:    aws.String(objKey),
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *S3Client) GetPresignURL(objKey string) (string, error) {
	presignClient := s3.NewPresignClient(s.s3)
	presignUrl, err := presignClient.PresignGetObject(context.Background(),
		&s3.GetObjectInput{
			Bucket: aws.String(s.cfg.BucketName),
			Key:    aws.String(objKey)},
		s3.WithPresignExpires(time.Hour*24))
	if err != nil {
		return "", err
	}

	return presignUrl.URL, nil
}
