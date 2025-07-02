# Module 1: Building a Transaction Feature (Go Data Types & Structs)

Welcome to your first feature module\! The goal here is to build a `transactions` endpoint from the ground up. In the process, you'll get hands-on experience with some of the most fundamental concepts in Go: its strong type system, how to define data structures (`structs`), and how to interact with a database using the existing patterns in this codebase.

For those coming from Laravel, you'll find the overall architecture familiar (Handler/Controller -\> Service -\> Model), but you'll notice Go requires us to be much more explicit about the *shape* and *type* of our data, which is a key to writing robust and bug-free applications.

### Learning Goals

  * **Understand Go's integer types** (`int8`, `int32`, `int64`) and when to use them.
  * **Define GORM model `structs`** that map directly to database tables.
  * **Create a custom `enum` type** using Go's `iota` keyword.
  * **Handle nullable booleans** in the database using pointers (`*bool`).
  * **Follow the project's Clean Architecture** to add a new feature.

## Step 1: Database Migration - Defining Our Schema

First, we need to tell the database about our new tables: `transactions` and `transaction_details`. We do this with SQL migration files.

1.  **Create the Migration Files**: Open your terminal in the project root. If you have `migrate` installed locally, you can use `migrate create`. For this guide, we'll just create the files directly.

2.  **Create the `transactions` table**:
    Create a new file: `database/migrations/000004_create_transactions_table.up.sql`

    ```sql
    CREATE TABLE transactions (
        id CHAR(36) PRIMARY KEY,
        user_id CHAR(36) NOT NULL,
        invoice_code VARCHAR(20) NOT NULL UNIQUE,
        total BIGINT NOT NULL,
        is_paid BOOLEAN NOT NULL DEFAULT FALSE,
        note TEXT,
        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
    );
    ```

3.  **Create the `transaction_details` table**:
    Create a new file: `database/migrations/000005_create_transaction_details_table.up.sql`

    ```sql
    CREATE TABLE transaction_details (
        id CHAR(36) PRIMARY KEY,
        transaction_id CHAR(36) NOT NULL,
        product_id CHAR(36) NOT NULL,
        product_name VARCHAR(255) NOT NULL,
        category TINYINT NOT NULL, -- For our enum
        qty TINYINT NOT NULL,
        price INT NOT NULL,
        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        FOREIGN KEY (transaction_id) REFERENCES transactions(id) ON DELETE CASCADE
    );
    ```

    > **💡 Why `BIGINT` and `TINYINT`?**
    > We use `BIGINT` for the transaction `total` because it could become a very large number. We use `TINYINT` for `qty` and `category` because we know they will always be small numbers, which saves space in the database.

4.  **Create the `down` migrations** (for rollbacks):

      * `000004_create_transactions_table.down.sql`:
        ```sql
        DROP TABLE IF EXISTS transactions;
        ```
      * `000005_create_transaction_details_table.down.sql`:
        ```sql
        DROP TABLE IF EXISTS transaction_details;
        ```

5.  **Run the Migrations**:
    Open a new terminal and run the migration command to apply the new tables to your database.

    ```bash
    docker-compose run --rm app go run ./cmd/migrate/main.go up
    ```

## Step 2: Model Creation - The Shape of Our Data

Now, let's create the Go `structs` that will represent our new tables. We'll follow the pattern from `internal/model/user_model.go`.

1.  **Create the `ProductCategory` Enum**: First, let's define our product categories. This is a perfect use case for `iota`, which creates incrementing numbers for us automatically.
    Create a new file: `internal/model/product_category.go`

    ```go
    package model

    // ProductCategory defines the enum for product types.
    type ProductCategory uint8 // Use uint8 since it's a small, non-negative number

    const (
    	_ ProductCategory = iota // Start with 0, but we'll ignore it
    	Goods
    	Service
    	Subscription
    )
    ```

2.  **Create the Transaction Models**: Now create two new model files.

      * `internal/model/transaction_model.go`:

    <!-- end list -->

    ```go
    package model

    import (
    	"context"
    	"time"

    	"github.com/google/uuid"
    	"gorm.io/gorm"
    )

    // Transaction defines the main transaction model.
    type Transaction struct {
    	ID                 uuid.UUID `gorm:"type:char(36);primary_key"`
    	UserID             uuid.UUID `gorm:"type:char(36);not null"`
    	InvoiceCode        string    `gorm:"size:20;not null;unique"`
    	Total              int64     `gorm:"not null"`
    	IsPaid             *bool     `gorm:"not null;default:false" json:"is_paid"`
    	Note               string    `gorm:"type:text"`
    	CreatedAt          time.Time
    	UpdatedAt          time.Time

    	// Relationships
    	User               User                `gorm:"foreignKey:UserID"`
    	TransactionDetails []TransactionDetail `gorm:"foreignKey:TransactionID"`
    }

    // BeforeCreate is a GORM hook.
    func (t *Transaction) BeforeCreate(tx *gorm.DB) (err error) {
    	t.ID = uuid.New()
    	// Set default IsPaid to false if it's nil
    	if t.IsPaid == nil {
    		b := false
    		t.IsPaid = &b
    	}
    	return
    }

    // Save creates or updates a record.
    func (t *Transaction) Save(db *gorm.DB) error {
    	return db.WithContext(context.Background()).Save(t).Error
    }
    ```

    > **💡 Why `*bool`?**
    > A pointer can be `nil`. A regular `bool` cannot. GORM often ignores fields with a "zero value" (like `false` for `bool`) during updates. Using a pointer (`*bool`) makes our intent clear: `nil` means "don't update," while a pointer to `false` means "update the value to `false`."

      * `internal/model/transaction_detail_model.go`:

    <!-- end list -->

    ```go
    package model

    import (
    	"time"
    	"github.com/google/uuid"
    	"gorm.io/gorm"
    )

    // TransactionDetail defines the items within a transaction.
    type TransactionDetail struct {
    	ID            uuid.UUID `gorm:"type:char(36);primary_key"`
    	TransactionID uuid.UUID `gorm:"type:char(36);not null"`
    	ProductID     uuid.UUID `gorm:"type:char(36);not null"`
    	ProductName   string    `gorm:"size:255;not null"`
    	Category      ProductCategory `gorm:"not null"`
    	Qty           int8      `gorm:"not null"`
    	Price         int32     `gorm:"not null"`
    	CreatedAt     time.Time
    	UpdatedAt     time.Time
    }

    // BeforeCreate is a GORM hook.
    func (td *TransactionDetail) BeforeCreate(tx *gorm.DB) (err error) {
    	td.ID = uuid.New()
    	return
    }
    ```

## Step 3: Service Layer - Implementing Business Logic

The service layer contains the core logic of our feature. It will calculate totals, generate invoice codes, and save everything to the database.

Create a new file: `internal/service/transaction_service.go`

```go
package service

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"
	"venturo-core/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TransactionService handles the business logic for transactions.
type TransactionService struct {
	db *gorm.DB
}

// NewTransactionService creates a new transaction service.
func NewTransactionService(db *gorm.DB) *TransactionService {
	return &TransactionService{db: db}
}

// CreateTransactionInput is the data structure needed to create a transaction.
type CreateTransactionInput struct {
	UserID uuid.UUID
	Items  []struct {
		ProductID   uuid.UUID
		ProductName string
		Category    model.ProductCategory
		Qty         int8
		Price       int32
	}
	Note string
}

// CreateTransaction handles the logic of creating a full transaction.
func (s *TransactionService) CreateTransaction(ctx context.Context, input CreateTransactionInput) (*model.Transaction, error) {
	// 1. Calculate the grand total and prepare details
	var total int64
	var details []model.TransactionDetail
	var itemNames []string

	for _, item := range input.Items {
		total += int64(item.Qty) * int64(item.Price)
		details = append(details, model.TransactionDetail{
			ProductID:   item.ProductID,
			ProductName: item.ProductName,
			Category:    item.Category,
			Qty:         item.Qty,
			Price:       item.Price,
		})
		itemNames = append(itemNames, item.ProductName)
	}

	// 2. Generate Invoice Code and Note
	invoiceCode := generateInvoiceCode()
	note := fmt.Sprintf("INV %s includes: %s. Additional notes: %s",
		invoiceCode,
		strings.Join(itemNames, ", "),
		input.Note,
	)

	// 3. Create the main transaction object
	transaction := model.Transaction{
		UserID:             input.UserID,
		InvoiceCode:        invoiceCode,
		Total:              total,
		Note:               note,
		TransactionDetails: details, // GORM will auto-create these
	}

	// 4. Save everything in a database transaction
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&transaction).Error; err != nil {
			return err
		}
		// The details should be created automatically due to the relationship,
		// but explicit creation is safer if auto-creation is disabled.
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &transaction, nil
}

// generateInvoiceCode creates a random invoice code.
func generateInvoiceCode() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("INV-%d-%04d", time.Now().Year(), rand.Intn(10000))
}
```

## Step 4: Handler - Exposing the Feature via API

The handler is the bridge between the web request and our service logic.

Create a new file: `internal/handler/http/transaction_handler.go`

```go
package http

import (
	"errors"
	"venturo-core/internal/model"
	"venturo-core/internal/service"
	"venturo-core/pkg/response"
	"venturo-core/pkg/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type TransactionHandler struct {
	transactionService *service.TransactionService
}

// NewTransactionHandler creates a new handler.
func NewTransactionHandler(s *service.TransactionService) *TransactionHandler {
	return &TransactionHandler{transactionService: s}
}

// CreateTransactionPayload defines the expected JSON from the client.
type CreateTransactionPayload struct {
	Items []struct {
		ProductID   uuid.UUID `json:"product_id" validate:"required"`
		ProductName string    `json:"product_name" validate:"required"`
		Category    uint8     `json:"category" validate:"required,min=1,max=3"`
		Qty         int8      `json:"qty" validate:"required,min=1"`
		Price       int32     `json:"price" validate:"required,min=0"`
	} `json:"items" validate:"required,min=1"`
	Note string `json:"note"`
}

// CreateTransaction is the handler for creating a new transaction.
// @Summary      Create a new transaction
// @Description  Creates a transaction with multiple detail items for the authenticated user.
// @Tags         Transactions
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        payload  body      CreateTransactionPayload  true  "Transaction Payload"
// @Success      201      {object}  response.ApiResponse{data=model.Transaction} "Successfully created transaction"
// @Failure      400      {object}  response.ApiResponse "Bad Request"
// @Failure      401      {object}  response.ApiResponse "Unauthorized"
// @Router       /transactions [post]
func (h *TransactionHandler) CreateTransaction(c *fiber.Ctx) error {
	userID, ok := c.Locals("current_user_id").(uuid.UUID)
	if !ok {
		return response.Error(c, fiber.StatusUnauthorized, errors.New("unauthorized"))
	}

	payload := new(CreateTransactionPayload)
	if err := c.BodyParser(payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, errors.New("cannot parse JSON"))
	}

	if errs := validator.ValidateStruct(payload); errs != nil {
		return response.ValidationError(c, errs)
	}

	// Map payload to the service input struct
	serviceInput := service.CreateTransactionInput{
		UserID: userID,
		Note:   payload.Note,
	}
	for _, item := range payload.Items {
		serviceInput.Items = append(serviceInput.Items, struct {
			ProductID   uuid.UUID
			ProductName string
			Category    model.ProductCategory
			Qty         int8
			Price       int32
		}{
			ProductID:   item.ProductID,
			ProductName: item.ProductName,
			Category:    model.ProductCategory(item.Category),
			Qty:         item.Qty,
			Price:       item.Price,
		})
	}

	// Call the service
	transaction, err := h.transactionService.CreateTransaction(c.Context(), serviceInput)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err)
	}

	return response.Success(c, fiber.StatusCreated, transaction)
}
```

## Step 5: Routing - Making It Accessible

The final step is to register our new endpoint so the application knows about it.

Open `internal/server/routes.go` and add the new service, handler, and route.

```go
package server

import (
	"sync"
	"venturo-core/configs"
	"venturo-core/internal/handler/http"
	"venturo-core/internal/middleware"
	"venturo-core/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"gorm.io/gorm"
)

func registerRoutes(app *fiber.App, db *gorm.DB, conf *configs.Config, wg *sync.WaitGroup) {
	// ... (keep all existing code)

	// --- Setup services ---
	authService := service.NewAuthService(db, conf)
	userService := service.NewUserService(db, wg)
	postService := service.NewPostService(db)
	// Add our new service
	transactionService := service.NewTransactionService(db) // <-- ADD THIS

	// --- Setup handlers ---
	authHandler := http.NewAuthHandler(authService)
	userHandler := http.NewUserHandler(userService)
	postHandler := http.NewPostHandler(postService)
	// Add our new handler
	transactionHandler := http.NewTransactionHandler(transactionService) // <-- ADD THIS

	// ... (keep existing auth, user, and post routes)

	// --- Transaction routes ---
	api.Post("/transactions", authMiddleware, transactionHandler.CreateTransaction) // <-- ADD THIS
}
```

## Conclusion

That's it\! You've successfully added a complete feature to the application. If you restart your server (`air` should do this automatically), you can now send a `POST` request to `/api/v1/transactions` with a valid JWT and a JSON body to create a new transaction.

You have learned:

  * How to define database schemas and map them to Go structs.
  * The difference between integer types and how to use them effectively.
  * How to create custom enum types with `iota`.
  * How to handle nullable database fields with pointers.
  * The flow of a request through the Handler, Service, and Model layers.

In the next module, we'll build on this by creating a reporting feature and introducing Go's powerful concurrency with **goroutines**.