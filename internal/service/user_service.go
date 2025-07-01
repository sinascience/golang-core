package service

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"sync"
	"venturo-core/internal/adapter/storage"
	"venturo-core/internal/model"
	"venturo-core/pkg/uploader"

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

func NewUserService(db *gorm.DB, wg *sync.WaitGroup) *UserService {
	gcsAdapter := storage.NewGCSAdapter("your-gcs-bucket-name")
	fileUploader := uploader.NewFileUploader(gcsAdapter, tempUploadPath)

	// Ensure the temporary upload directory exists
	if err := os.MkdirAll(tempUploadPath, os.ModePerm); err != nil {
		log.Fatalf("could not create temp upload directory: %v", err)
	}
	return &UserService{db: db, uploader: fileUploader, wg: wg}
}

// GetUserProfile retrieves a user's profile by their ID.
func (s *UserService) GetUserProfile(userID uuid.UUID) (*model.User, error) {
	var user model.User
	return user.FindByID(s.db, userID)
}

// UpdateUserProfile updates a user's profile data.
func (s *UserService) UpdateUserProfile(ctx context.Context, userID uuid.UUID, newName string, file *multipart.FileHeader) (*model.User, error) {
	// First, find the user to ensure they exist.
	user, err := s.GetUserProfile(userID)
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
	if err := user.Save(s.db); err != nil {
		return nil, err
	}

	if file != nil {
		// Define the database logic in callbacks
		onLocalUpload := func() {
			user.ImageStatus = "local"
			if err := user.Save(s.db); err != nil {
				log.Printf("Error updating status to 'local' for user %s: %v", userID, err)
			}
		}

		onCloudUpload := func() {
			user.ImageStatus = "cloud"
			if err := user.Save(s.db); err != nil {
				log.Printf("Error updating status to 'cloud' for user %s: %v", userID, err)
			}
		}

		// Start the background process
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			s.uploader.UploadAsync(file, user.AvatarURL, onLocalUpload, onCloudUpload)
		}()
	}

	return user, nil
}
