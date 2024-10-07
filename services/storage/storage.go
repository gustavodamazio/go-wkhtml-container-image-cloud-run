package storage

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/storage"
)

// StorageService defines the structure for the Google Cloud Storage service
type StorageService struct {
	client *storage.Client
	bucket *storage.BucketHandle
	ctx    context.Context
}

// NewStorageService initializes and returns a new instance of StorageService
func NewStorageService(ctx context.Context, bucketName string) (*StorageService, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create storage client: %v", err)
	}
	bucket := client.Bucket(bucketName)

	return &StorageService{
		client: client,
		bucket: bucket,
		ctx:    ctx,
	}, nil
}

// ReadFile reads a file from Google Cloud Storage and returns its content
func (s *StorageService) ReadFile(fileName string) ([]byte, error) {
	reader, err := s.bucket.Object(fileName).NewReader(s.ctx)
	if err != nil {
		return nil, fmt.Errorf("could not open file from bucket: %v", err)
	}
	defer reader.Close()

	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("could not read file: %v", err)
	}

	return content, nil
}

// WriteFile writes the provided content to a file in Google Cloud Storage
func (s *StorageService) WriteFile(fileName string, content []byte) error {
	writer := s.bucket.Object(fileName).NewWriter(s.ctx)

	_, err := writer.Write(content)
	if err != nil {
		return fmt.Errorf("could not write file: %v", err)
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("could not close writer: %v", err)
	}

	return nil
}

// Close closes the storage client connection
func (s *StorageService) Close() error {
	return s.client.Close()
}
