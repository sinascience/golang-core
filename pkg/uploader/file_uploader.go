package uploader

import (
	"context"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"venturo-core/internal/adapter/storage"
)

// FileUploader handles the low-level logic of saving and uploading files.
type FileUploader struct {
	storageAdapter storage.StorageAdapter
	localPath      string
}

// NewFileUploader creates a new uploader instance.
func NewFileUploader(storageAdapter storage.StorageAdapter, localPath string) *FileUploader {
	// Ensure the temporary upload directory exists
	if err := os.MkdirAll(localPath, os.ModePerm); err != nil {
		log.Fatalf("could not create temp upload directory: %v", err)
	}
	return &FileUploader{storageAdapter: storageAdapter, localPath: localPath}
}

// UploadAsync handles the entire file processing pipeline in the background.
func (u *FileUploader) UploadAsync(
	file *multipart.FileHeader,
	objectName string,
	onLocalUploadSuccess func(),
	onCloudUploadSuccess func(),
) {
	// --- Step 1: Save Locally ---
	localFilePath := filepath.Join(u.localPath, objectName)
	if err := saveToLocal(file, localFilePath); err != nil {
		log.Printf("Error saving temp file %s: %v", objectName, err)
		return
	}
	log.Printf("Successfully saved temp file: %s", localFilePath)

	// Execute the first callback
	onLocalUploadSuccess()

	// --- Step 2: Upload to Cloud ---
	// In a real app, you would read the local file's contents to upload.
	// We pass the header for placeholder simplicity.
	if _, err := u.storageAdapter.Upload(context.Background(), file); err != nil {
		log.Printf("Error uploading to cloud for %s: %v", objectName, err)
		return // Don't continue if cloud upload fails
	}
	log.Printf("Successfully uploaded to cloud: %s", objectName)

	// Execute the second callback
	onCloudUploadSuccess()

	// --- Step 3: Cleanup ---
	if err := os.Remove(localFilePath); err != nil {
		log.Printf("Error cleaning up temp file %s: %v", objectName, err)
	}
	log.Printf("Successfully cleaned up temp file: %s", objectName)
}

// saveToLocal is a helper function containing the file-saving logic.
func saveToLocal(file *multipart.FileHeader, path string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(path)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	return err
}
