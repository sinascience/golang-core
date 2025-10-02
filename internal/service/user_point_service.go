package service

import (
	// "errors"
	"venturo-core/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserPointService struct {
	db *gorm.DB
}

// NewUserPointService creates a new user point service.
func NewUserPointService(db *gorm.DB) *UserPointService {
	return &UserPointService{db: db}
}

// CreateUserPoint creates a new user point record for a given user.
func (s *UserPointService) CreateUserPoint(userID uuid.UUID) (*model.UserPoint, error) {
	userPoint := model.UserPoint{UserID: userID, Point: 0}
	if err := userPoint.Save(s.db); err != nil {
		return nil, err
	}
	return &userPoint, nil
}

// FindPointByUserID retrieves a single user point by their userID.
func (s *UserPointService) FindPointByUserID(userID uuid.UUID) (*model.UserPoint, error) {
	var userPoint model.UserPoint
	return userPoint.FindByUserID(s.db, userID)
}

func (s *UserPointService) UpdateUserPoint(userID, pointID uuid.UUID, points int32) (*model.UserPoint, error) {
	// Find the user point first
	userPoint, err := s.FindPointByUserID(userID)
	if err != nil {
		return nil, err
	}

	// Update the fileds
	userPoint.Point = points

	// Simpan perubahan tanpa mengubah user_id
	if err := s.db.Model(&userPoint).Select("Point").Updates(userPoint).Error; err != nil {
		return nil, err
	}

	return userPoint, nil
}