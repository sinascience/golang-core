package model

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Post defines the post model.
type Post struct {
	ID        uuid.UUID `gorm:"type:char(36);primary_key" json:"id"`
	Title     string    `gorm:"size:255;not null" json:"title"`
	Body      string    `gorm:"type:text" json:"body"`
	UserID    uuid.UUID `gorm:"type:char(36);not null" json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Define the relationship to the User model
	User User `gorm:"foreignKey:UserID" json:"author,omitempty"`
}

// BeforeCreate is a GORM hook that runs before a new record is created.
func (p *Post) BeforeCreate(tx *gorm.DB) (err error) {
	p.ID = uuid.New()
	return
}

// Save creates or updates a post record.
func (p *Post) Save(db *gorm.DB) error {
	return db.WithContext(context.Background()).Save(p).Error
}

// FindAll retrieves all post records, preloading the author data.
func (p *Post) FindAll(db *gorm.DB, page, limit int) ([]Post, int64, error) {
	var posts []Post
	var total int64

	// 1. Get the total count of posts
	if err := db.Model(&Post{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 2. Calculate the offset for pagination
	offset := (page - 1) * limit

	// 3. Get the paginated data
	err := db.Limit(limit).Offset(offset).Preload("User").Order("created_at desc").Find(&posts).Error
	if err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

// FindByID retrieves a single post by its ID, preloading the author.
func (p *Post) FindByID(db *gorm.DB, id uuid.UUID) (*Post, error) {
	var post Post
	err := db.Preload("User").Where("id = ?", id).First(&post).Error
	return &post, err
}

// Delete removes a post record from the database.
func (p *Post) Delete(db *gorm.DB) error {
	return db.Delete(p).Error
}
