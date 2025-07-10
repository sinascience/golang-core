package storage

import (
	"context"
	"log/slog"
	"mime/multipart"
)

// StorageAdapter defines the interface for any cloud storage service.
type StorageAdapter interface {
	Upload(ctx context.Context, file *multipart.FileHeader, objectName string) (string, error)
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
func (a *GCSAdapter) Upload(ctx context.Context, file *multipart.FileHeader, objectName string) (string, error) {
	slog.Info("Simulating upload of file to GCS bucket", "file", file.Filename, "bucket", a.bucketName)
	dummyURL := "https://storage.googleapis.com/" + a.bucketName + "/" + file.Filename
	return dummyURL, nil
}
