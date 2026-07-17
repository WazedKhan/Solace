package storage

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Storage struct {
	client     *s3.Client
	presigner  *s3.PresignClient
	bucketName string
}

func NewS3Storage(
	client *s3.Client,
	presigner *s3.PresignClient,
	bucketName string,
) *S3Storage {
	return &S3Storage{
		client:     client,
		presigner:  presigner,
		bucketName: bucketName,
	}
}
