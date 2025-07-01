package storage

import (
	"context"
	"log"
	"mime/multipart"
)

// StorageAdapter defines the interface for any cloud storage service.
type StorageAdapter interface {
	Upload(ctx context.Context, file *multipart.FileHeader) (string, error)
}

// GCSAdapter is the placeholder implementation for Google Cloud Storage.
type GCSAdapter struct {
	bucketName string
}

// NewGCSAdapter creates a new instance of our GCS placeholder.
func NewGCSAdapter(bucketName string) *GCSAdapter {
	return &GCSAdapter{bucketName: bucketName}
}

// Upload simulates uploading a file to GCS.
func (a *GCSAdapter) Upload(ctx context.Context, file *multipart.FileHeader) (string, error) {
	// For now, this is a placeholder. It just logs the action.
	// In a real implementation, this would contain the GCS upload logic.
	log.Printf("Simulating upload of file '%s' to GCS bucket '%s'.", file.Filename, a.bucketName)

	// We return a dummy URL and a nil error to simulate success.
	dummyURL := "https://storage.googleapis.com/" + a.bucketName + "/" + file.Filename
	return dummyURL, nil
}
