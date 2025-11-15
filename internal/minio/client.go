
package minio

import (
	"context"
	"log"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var Client *minio.Client

// NewClient creates and returns a new MinIO client
func NewClient(endpoint, accessKeyID, secretAccessKey string, useSSL bool) {
	var err error
	Client, err = minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}
}

// GetPresignedURL returns a presigned URL for the given object
func GetPresignedURL(bucketName, objectName string) (string, error) {
	presignedURL, err := Client.PresignedGetObject(context.Background(), bucketName, objectName, time.Second*24*60*60, nil)
	if err != nil {
		return "", err
	}
	return presignedURL.String(), nil
}
