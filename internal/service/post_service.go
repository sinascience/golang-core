package service

import (
	"context"
	"errors"
	"venturo-core/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostService struct {
	db *gorm.DB
}

// NewPostService creates a new post service.
func NewPostService(db *gorm.DB) *PostService {
	return &PostService{db: db}
}

// CreatePost creates a new post for a given user.
func (s *PostService) CreatePost(ctx context.Context, userID uuid.UUID, title, body string) (*model.Post, error) {
	post := model.Post{
		Title:  title,
		Body:   body,
		UserID: userID,
	}

	if err := post.Save(ctx, s.db); err != nil {
		return nil, err
	}
	return &post, nil
}

// GetAllPosts retrieves all posts.
func (s *PostService) GetAllPosts(ctx context.Context, page, limit int) ([]model.Post, int64, error) {
	var post model.Post
	return post.FindAll(ctx, s.db, page, limit)
}

// GetPostByID retrieves a single post by its ID.
func (s *PostService) GetPostByID(ctx context.Context, id uuid.UUID) (*model.Post, error) {
	var post model.Post
	return post.FindByID(ctx, s.db, id)
}

// DeletePost finds a post, checks for ownership, and deletes it.
func (s *PostService) DeletePost(ctx context.Context, postID, userID uuid.UUID) error {
	// Find the post first
	post, err := s.GetPostByID(ctx, postID)
	if err != nil {
		return err // Post not found
	}

	// Authorization Check: Ensure the user owns the post
	if post.UserID != userID {
		return errors.New("unauthorized: you are not the owner of this post")
	}

	// Delete the post
	return post.Delete(ctx, s.db)
}

// UpdatePost finds a post, checks for ownership, and updates it.
func (s *PostService) UpdatePost(ctx context.Context, postID, userID uuid.UUID, newTitle, newBody string) (*model.Post, error) {
	// Find the post first
	post, err := s.GetPostByID(ctx, postID)
	if err != nil {
		return nil, err // Post not found
	}

	// Authorization Check: Ensure the user owns the post
	if post.UserID != userID {
		return nil, errors.New("unauthorized: you are not the owner of this post")
	}

	// Update only the title and body fields using GORM's Updates method
	// This prevents the user_id from being modified
	if err := s.db.WithContext(ctx).Model(post).Updates(model.Post{
		Title: newTitle,
		Body:  newBody,
	}).Error; err != nil {
		return nil, err
	}

	// Update the local object to reflect the changes
	post.Title = newTitle
	post.Body = newBody

	return post, nil
}
