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
// @Summary      Create a new post
// @Description  Creates a new post for the authenticated user.
// @Tags         Posts
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        payload  body      CreatePostPayload      true  "Post Creation Payload"
// @Success      201      {object}  response.ApiResponse{data=model.Post} "Successfully created post"
// @Failure      400      {object}  response.ApiResponse "Bad Request"
// @Failure      401      {object}  response.ApiResponse "Unauthorized"
// @Router       /posts [post]
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
// @Summary      Get all posts
// @Description  Retrieves a paginated list of all posts.
// @Tags         Posts
// @Produce      json
// @Param        page   query     int  false  "Page number for pagination" default(1)
// @Param        limit  query     int  false  "Number of items per page" default(10)
// @Success      200    {object}  response.ApiResponse{data=[]model.Post} "Successfully retrieved posts"
// @Failure      500    {object}  response.ApiResponse "Internal Server Error"
// @Router       /posts [get]
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
// @Summary      Get a single post
// @Description  Retrieves a single post by its unique ID.
// @Tags         Posts
// @Produce      json
// @Param        id   path      string  true  "Post ID"
// @Success      200  {object}  response.ApiResponse{data=model.Post} "Successfully retrieved post"
// @Failure      404  {object}  response.ApiResponse "Post not found"
// @Router       /posts/{id} [get]
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
// @Summary      Delete a post
// @Description  Deletes a post. Only the author can delete their post.
// @Tags         Posts
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id   path      string  true  "Post ID"
// @Success      200  {object}  response.ApiResponse "Successfully deleted post"
// @Failure      401      {object}  response.ApiResponse "Unauthorized"
// @Failure      403      {object}  response.ApiResponse "Forbidden"
// @Failure      404      {object}  response.ApiResponse "Post not found"
// @Router       /posts/{id} [delete]
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
// @Summary      Update a post
// @Description  Updates a post. Only the author can update their post.
// @Tags         Posts
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id       path      string               true  "Post ID"
// @Param        payload  body      CreatePostPayload    true  "Post Update Payload"
// @Success      200      {object}  response.ApiResponse{data=model.Post} "Successfully updated post"
// @Failure      401      {object}  response.ApiResponse "Unauthorized"
// @Failure      403      {object}  response.ApiResponse "Forbidden"
// @Failure      404      {object}  response.ApiResponse "Post not found"
// @Router       /posts/{id} [put]
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
