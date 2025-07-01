package http

import (
	"errors"
	"strconv"
	"strings"
	"venturo-core/internal/model"
	"venturo-core/internal/service"
	"venturo-core/pkg/response"
	"venturo-core/pkg/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// PaginatedPostsResponse defines the structure for a paginated API response.
type PaginatedPostsResponse struct {
	Data []model.Post `json:"data"`
	Meta Meta         `json:"meta"`
}

// Meta holds the pagination metadata.
type Meta struct {
	TotalRecords int64 `json:"total_records"`
	CurrentPage  int   `json:"current_page"`
	PerPage      int   `json:"per_page"`
	TotalPages   int   `json:"total_pages"`
}

type PostHandler struct {
	postService *service.PostService
}

// NewPostHandler creates a new PostHandler.
func NewPostHandler(postService *service.PostService) *PostHandler {
	return &PostHandler{postService: postService}
}

// CreatePostPayload defines the expected JSON for creating a post.
type CreatePostPayload struct {
	Title string `json:"title" validate:"required,min=5"`
	Body  string `json:"body"`
}

// CreatePost is the handler for creating a new post.
func (h *PostHandler) CreatePost(c *fiber.Ctx) error {
	// Get user ID from the JWT middleware
	userID, ok := c.Locals("current_user_id").(uuid.UUID)
	if !ok {
		return response.Error(c, fiber.StatusUnauthorized, errors.New("unauthorized"))
	}

	payload := new(CreatePostPayload)
	if err := c.BodyParser(payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, errors.New("cannot parse JSON"))
	}

	if errs := validator.ValidateStruct(payload); errs != nil {
		return response.ValidationError(c, errs)
	}

	post, err := h.postService.CreatePost(userID, payload.Title, payload.Body)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, errors.New("could not create post"))
	}

	return response.Success(c, fiber.StatusCreated, post)
}

// GetAllPosts now handles pagination and returns a structured response.
func (h *PostHandler) GetAllPosts(c *fiber.Ctx) error {
	// 1. Parse query parameters for pagination
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit < 1 {
		limit = 10
	}
	if limit > 100 { // Set a max limit
		limit = 100
	}

	// 2. Call the service to get paginated data and total count
	posts, total, err := h.postService.GetAllPosts(page, limit)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, errors.New("could not retrieve posts"))
	}

	return response.Pagination(c, posts, page, limit, total)
}

// GetPostByID is the handler for retrieving a single post by its ID.
func (h *PostHandler) GetPostByID(c *fiber.Ctx) error {
	// Get ID from URL parameter
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, errors.New("invalid ID format"))
	}

	post, err := h.postService.GetPostByID(id)
	if err != nil {
		return response.Error(c, fiber.StatusNotFound, errors.New("post not found"))
	}

	return response.Success(c, fiber.StatusOK, post)
}

// DeletePost is the handler for deleting a post.
func (h *PostHandler) DeletePost(c *fiber.Ctx) error {
	// Get post ID from URL parameter
	postID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, errors.New("invalid ID format"))
	}

	// Get user ID from the JWT middleware
	userID, ok := c.Locals("current_user_id").(uuid.UUID)
	if !ok {
		return response.Error(c, fiber.StatusUnauthorized, errors.New("unauthorized"))
	}

	// Call the service to delete the post
	err = h.postService.DeletePost(postID, userID)
	if err != nil {
		if strings.Contains(err.Error(), "unauthorized") {
			return response.Error(c, fiber.StatusForbidden, err)
		}
		if strings.Contains(err.Error(), "not found") {
			return response.Error(c, fiber.StatusNotFound, errors.New("post not found"))
		}
		return response.Error(c, fiber.StatusInternalServerError, errors.New("could not delete post"))
	}

	// Return 204 No Content for successful deletion
	return response.Success(c, fiber.StatusOK, nil)
}

// UpdatePost is the handler for updating a post.
func (h *PostHandler) UpdatePost(c *fiber.Ctx) error {
	// Get post ID from URL parameter
	postID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, errors.New("invalid ID format"))
	}

	// Get user ID from the JWT middleware
	userID, ok := c.Locals("current_user_id").(uuid.UUID)
	if !ok {
		return response.Error(c, fiber.StatusUnauthorized, errors.New("unauthorized"))
	}

	payload := new(CreatePostPayload)
	if err := c.BodyParser(payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, errors.New("cannot parse JSON"))
	}

	if errs := validator.ValidateStruct(payload); errs != nil {
		return response.ValidationError(c, errs)
	}

	updatedPost, err := h.postService.UpdatePost(postID, userID, payload.Title, payload.Body)
	if err != nil {
		if strings.Contains(err.Error(), "unauthorized") {
			return response.Error(c, fiber.StatusForbidden, err)
		}
		if strings.Contains(err.Error(), "not found") {
			return response.Error(c, fiber.StatusNotFound, errors.New("post not found"))
		}
		return response.Error(c, fiber.StatusInternalServerError, errors.New("could not update post"))
	}

	return response.Success(c, fiber.StatusOK, updatedPost)
}
