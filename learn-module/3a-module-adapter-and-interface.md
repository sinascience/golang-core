# Module 3: Adapters, Concurrency, and State

Welcome to the final module\! This is where all the concepts you've learned come together. We will build a feature to create a new product with an image upload. This single feature will use:

  * **The Adapter Pattern**: To decouple our application from the file storage system.
  * **Concurrency**: To upload the image in the background using a **goroutine** so the user gets an instant response.
  * **State Management**: To track the status of the image upload directly in our database model, just like the user's avatar.

This is how professional, robust, and scalable features are built.

### Learning Goals

  * Implement the **Adapter Pattern** for a third-party service.
  * Combine the adapter with a **goroutine** to handle uploads asynchronously.
  * Manage the **state** of a model (`ImageStatus`) throughout an asynchronous process.
  * Use `sync.WaitGroup` to ensure background uploads can complete gracefully.
  * Use **Dependency Injection** to provide multiple dependencies to a service.

## Step 1: Database Migration and Model with State

First, let's create our `products` table. This time, we will include the `image_status` field to track the upload process.

1.  **Create the Migration File**:
    Create a new file: `database/migrations/000007_create_products_table.up.sql`

    ```sql
    CREATE TABLE products (
        id CHAR(36) PRIMARY KEY,
        name VARCHAR(255) NOT NULL,
        price INT NOT NULL DEFAULT 0,
        stock SMALLINT NOT NULL DEFAULT 0,
        image_url VARCHAR(255),
        image_status VARCHAR(20) NOT NULL DEFAULT 'default',
        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
    );
    ```

2.  **Create the `down` migration**:
    Create `database/migrations/000007_create_products_table.down.sql`:

    ```sql
    DROP TABLE IF EXISTS products;
    ```

3.  **Run the Migration**:

    ```bash
    docker-compose run --rm app go run ./cmd/migrate/main.go up
    ```

4.  **Create the Product Model**:
    Create a new file: `internal/model/product_model.go`

    ```go
    package model

    import (
    	"context"
    	"time"

    	"github.com/google/uuid"
    	"gorm.io/gorm"
    )

    // Product defines the product model with image status.
    type Product struct {
    	ID           uuid.UUID `gorm:"type:char(36);primary_key"`
    	Name         string    `gorm:"size:255;not null"`
    	Price        int32
    	Stock        int16
    	ImageURL     string    `gorm:"size:255"`
    	ImageStatus  string    `gorm:"size:20;not null;default:'default'"`
    	CreatedAt    time.Time
    	UpdatedAt    time.Time
    }

    // BeforeCreate is a GORM hook.
    func (p *Product) BeforeCreate(tx *gorm.DB) (err error) {
    	p.ID = uuid.New()
    	return
    }

    // Save creates or updates a product record.
    func (p *Product) Save(db *gorm.DB) error {
    	return db.WithContext(context.Background()).Save(p).Error
    }
    ```

## Step 2: The Adapter

We will create a `LocalUploaderAdapter` that implements the existing `storage.StorageAdapter` interface.

Create a new file: `internal/adapter/storage/local_adapter.go`

```go
package storage

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"os"
	"path/filepath"
)

// LocalUploaderAdapter saves files to the local disk.
type LocalUploaderAdapter struct {
	basePath string
}

// NewLocalUploaderAdapter creates a new local uploader.
func NewLocalUploaderAdapter(basePath string) *LocalUploaderAdapter {
	if err := os.MkdirAll(basePath, os.ModePerm); err != nil {
		slog.Error("could not create local upload directory", "error", err)
		os.Exit(1)
	}
	return &LocalUploaderAdapter{basePath: basePath}
}

// Upload implements the StorageAdapter interface.
func (a *LocalUploaderAdapter) Upload(ctx context.Context, file *multipart.FileHeader, objectName string) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	dstPath := filepath.Join(a.basePath, objectName)
	dst, err := os.Create(dstPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return "", err
	}

	slog.Info("Successfully saved file to local disk", "path", dstPath)
	publicURL := "/" + filepath.ToSlash(dstPath)
	return publicURL, nil
}
```

> **Note**: As before, ensure your `StorageAdapter` interface in `internal/adapter/storage/storage_adapter.go` accepts `objectName` so the method signatures match.

## Step 3: Service Layer - Combining All Concepts

Now we create a `ProductService` that brings everything together. It will depend on the database, the WaitGroup, and the storage interface.

Create a new file: `internal/service/product_service.go`

```go
package service

import (
	"context"
	"fmt"
	"log/slog"
	"mime/multipart"
	"path/filepath"
	"sync"
	"venturo-core/internal/adapter/storage"
	"venturo-core/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ProductService handles the business logic for products.
type ProductService struct {
	db             *gorm.DB
	wg             *sync.WaitGroup
	storageAdapter storage.StorageAdapter
}

// NewProductService creates a new product service.
func NewProductService(db *gorm.DB, wg *sync.WaitGroup, storageAdapter storage.StorageAdapter) *ProductService {
	return &ProductService{db: db, wg: wg, storageAdapter: storageAdapter}
}

// CreateProductInput is the data needed to create a new product.
type CreateProductInput struct {
	Name  string
	Price int32
	Stock int16
	Image *multipart.FileHeader
}

// CreateProduct creates a product and asynchronously uploads its image.
func (s *ProductService) CreateProduct(ctx context.Context, input CreateProductInput) (*model.Product, error) {
	product := model.Product{
		Name:  input.Name,
		Price: input.Price,
		Stock: input.Stock,
	}

	// If an image is provided, prepare for upload.
	if input.Image != nil {
		imageName := fmt.Sprintf("%s%s", uuid.NewString(), filepath.Ext(input.Image.Filename))
		product.ImageURL = imageName
		product.ImageStatus = "uploading"
	}

	// Save the initial product record. This is fast and synchronous.
	if err := product.Save(s.db); err != nil {
		return nil, err
	}

	// If there's an image, start the background upload process.
	if input.Image != nil {
		s.wg.Add(1)
		go s.uploadProductImage(product.ID, input.Image, product.ImageURL)
	}

	return &product, nil
}

// uploadProductImage is the background worker.
func (s *ProductService) uploadProductImage(productID uuid.UUID, file *multipart.FileHeader, objectName string) {
	defer s.wg.Done()
	bgCtx := context.Background()

	// 1. Upload the file.
	publicURL, err := s.storageAdapter.Upload(bgCtx, file, objectName)
	if err != nil {
		slog.Error("Failed to upload product image", "productID", productID, "error", err)
		// Update status to 'failed'
		s.updateImageStatus(productID, "failed", "")
		return
	}

	// 2. Update status to 'done' on success.
	slog.Info("Successfully uploaded product image", "productID", productID, "url", publicURL)
	s.updateImageStatus(productID, "done", publicURL)
}

// updateImageStatus is a helper to update the product record.
func (s *ProductService) updateImageStatus(productID uuid.UUID, status string, url string) {
	var product model.Product
	// Using a map to update specific fields.
	updates := map[string]interface{}{"image_status": status}
	if url != "" {
		updates["image_url"] = url
	}

	if err := s.db.Model(&product).Where("id = ?", productID).Updates(updates).Error; err != nil {
		slog.Error("Failed to update product image status", "productID", productID, "error", err)
	}
}
```

## Step 4: Handler & Dependency Injection

Finally, we wire everything together in the handler and router.

1.  **Create the Product Handler**:
    Create a new file: `internal/handler/http/product_handler.go` (This is the same as the previous version).

    ```go
    package http

    import (
    	"errors"
    	"strconv"
    	"venturo-core/internal/service"
    	"venturo-core/pkg/response"

    	"github.com/gofiber/fiber/v2"
    )

    type ProductHandler struct {
    	productService *service.ProductService
    }

    func NewProductHandler(s *service.ProductService) *ProductHandler {
    	return &ProductHandler{productService: s}
    }

    // CreateProduct handles the multipart/form-data request to create a product.
    // ... (Add Swagger annotations here) ...
    func (h *ProductHandler) CreateProduct(c *fiber.Ctx) error {
    	price, err := strconv.Atoi(c.FormValue("price"))
    	if err != nil {
    		return response.Error(c, fiber.StatusBadRequest, errors.New("invalid price format"))
    	}
    	stock, err := strconv.Atoi(c.FormValue("stock"))
    	if err != nil {
    		return response.Error(c, fiber.StatusBadRequest, errors.New("invalid stock format"))
    	}

    	input := service.CreateProductInput{
    		Name:  c.FormValue("name"),
    		Price: int32(price),
    		Stock: int16(stock),
    	}

    	file, err := c.FormFile("image")
    	if err == nil {
    		input.Image = file
    	}

    	product, err := h.productService.CreateProduct(c.Context(), input)
    	if err != nil {
    		return response.Error(c, fiber.StatusInternalServerError, err)
    	}

    	return response.Success(c, fiber.StatusCreated, product)
    }
    ```

2.  **Wire It Up in `routes.go`**:
    Open `internal/server/routes.go` and inject all three dependencies into the `ProductService`.

    ```go
    // Add "venturo-core/internal/adapter/storage" to imports

    func registerRoutes(app *fiber.App, db *gorm.DB, conf *configs.Config, wg *sync.WaitGroup) {
    	// ... (keep all existing code)

    	// --- Setup Adapters ---
    	localUploader := storage.NewLocalUploaderAdapter("./public/uploads")

    	// --- Setup services ---
    	// ...
    	transactionService := service.NewTransactionService(db, wg)
    	// Inject all three dependencies into the ProductService.
    	productService := service.NewProductService(db, wg, localUploader) // <-- INJECT ALL 3

    	// --- Setup handlers ---
    	// ...
    	transactionHandler := http.NewTransactionHandler(transactionService)
    	productHandler := http.NewProductHandler(productService)

    	// ... (keep existing routes)

    	// --- Product routes ---
    	api.Post("/products", authMiddleware, productHandler.CreateProduct)
    }
    ```

## Conclusion

Congratulations\! You have completed the Venturo Golang Core learning journey. You have built a feature that demonstrates the most important patterns for modern backend development: a clean separation of concerns, robust state management, safe concurrency, and loose coupling from third-party services. You are now well-equipped to build professional, high-quality applications in Go.

-----