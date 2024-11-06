package backend

import (
	"context"
	"fmt"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var bucketName string = "code-inspector"

func InitializeMinIO() (*minio.Client, error) {
	endpoint, accessKey, secretKey, useSSL, err := FromEnvMinIO()
	if err != nil {
		return nil, fmt.Errorf("failed to get MinIO configuration from environment: %v", err)
	}

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize MinIO client: %v", err)
	}

	log.Println("MinIO client initialized successfully")

	found, err := client.BucketExists(context.Background(), bucketName)
	if err != nil {
		return nil, fmt.Errorf("failed to check if bucket exists: %v", err)
	}

	if !found {
		err = client.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to create bucket: %v", err)
		}
		log.Printf("Bucket '%s' created successfully\n", bucketName)
	} else {
		log.Printf("Bucket '%s' already exists\n", bucketName)
	}

	return client, nil
}
