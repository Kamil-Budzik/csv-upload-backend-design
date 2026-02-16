package storage

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/kamil-budzik/csv-processor/internal/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioStorage struct {
	bucketName string
	client     *minio.Client
}

func NewMinioStorage(client *minio.Client, bucketName string) *MinioStorage {
	if client == nil {
		panic("NewMinioStorage: Minio client cannot be nil")
	}
	if len(bucketName) == 0 {
		panic("NewMinioStorage: Bucket name cannot be empty")
	}

	return &MinioStorage{client: client, bucketName: bucketName}
}

func Connect(cfg config.Config) *minio.Client {
	endpoint := fmt.Sprintf("localhost:%s", cfg.MinioApiPort)
	useSSL := false

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinioUser, cfg.MinioPassword, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Minio Client is connected: %s\n", minioClient.EndpointURL())
	return minioClient
}

func (ms MinioStorage) CreateBucket(ctx context.Context) {
	location := "us-east-1"
	err := ms.client.MakeBucket(ctx, ms.bucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		exists, errBucketExists := ms.client.BucketExists(ctx, ms.bucketName)
		if errBucketExists == nil && exists {
			log.Println("Bucket already exists")
			return
		} else {
			log.Fatal(err)
			return
		}
	}
	log.Println("Bucket created successfully!")
}

func (ms MinioStorage) UploadCSV(ctx context.Context, fileName string, fileSize int64, file io.Reader) (string, error) {
	uploadInfo, err := ms.client.PutObject(ctx, ms.bucketName, fileName, file, fileSize, minio.PutObjectOptions{
		ContentType: "text/csv",
	})
	if err != nil {
		log.Printf("Error in Uploading File %v\n", err)
	}
	return uploadInfo.Location, err

}
