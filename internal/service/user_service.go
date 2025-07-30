package service

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"sync"
	"venturo-core/internal/adapter/storage"
	"venturo-core/internal/model"
	"venturo-core/pkg/uploader"

	"log/slog"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Define the temporary local storage path
const tempUploadPath = "./public/uploads/avatars"

type UserService struct {
	db       *gorm.DB
	uploader *uploader.FileUploader
	wg       *sync.WaitGroup
}

func NewUserService(db *gorm.DB, bucketName string, wg *sync.WaitGroup) *UserService {
	gcsAdapter := storage.NewGCSAdapter(bucketName)
	fileUploader := uploader.NewFileUploader(gcsAdapter, tempUploadPath)

	// Ensure the temporary upload directory exists
	if err := os.MkdirAll(tempUploadPath, os.ModePerm); err != nil {
		slog.Error("could not create temp upload directory", "error", err)
		os.Exit(1)
	}
	return &UserService{db: db, uploader: fileUploader, wg: wg}
}

// GetUserProfile retrieves a user's profile by their ID.
func (s *UserService) GetUserProfile(ctx context.Context, userID uuid.UUID) (*model.User, error) {
	var user model.User
	return user.FindByID(ctx, s.db, userID)
}

// UpdateUserProfile updates a user's profile data.
func (s *UserService) UpdateUserProfile(ctx context.Context, userID uuid.UUID, newName string, file *multipart.FileHeader) (*model.User, error) {
	// First, find the user to ensure they exist.
	user, err := s.GetUserProfile(ctx, userID)
	if err != nil {
		return nil, err // User not found
	}

	if file != nil {
		// Generate a new unique filename
		ext := filepath.Ext(file.Filename)
		newFileName := fmt.Sprintf("%s%s", uuid.New().String(), ext)
		user.AvatarURL = newFileName // Store only the filename
		user.ImageStatus = "uploading"
	}

	// Update the user's name.
	user.Name = newName

	// Save the updated user record. GORM's Save() handles updates automatically.
	if err := user.Save(ctx, s.db); err != nil {
		return nil, err
	}

	if file != nil {
		// Capture userID and filename for background process to avoid race conditions
		capturedUserID := userID
		capturedFilename := user.AvatarURL
		
		// Define the database logic in callbacks with proper error handling
		onLocalUpload := func() {
			if err := s.db.Model(&model.User{}).Where("id = ?", capturedUserID).Update("image_status", "local").Error; err != nil {
				slog.Error("Error updating status to 'local' for user", "userID", capturedUserID, "error", err)
			}
		}

		onCloudUpload := func() {
			if err := s.db.Model(&model.User{}).Where("id = ?", capturedUserID).Update("image_status", "cloud").Error; err != nil {
				slog.Error("Error updating status to 'cloud' for user", "userID", capturedUserID, "error", err)
			}
		}

		// Start the background process
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			s.uploader.UploadAsync(file, capturedFilename, onLocalUpload, onCloudUpload)
		}()
	}

	return user, nil
}
