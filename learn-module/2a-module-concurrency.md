# Module 2: Concurrency & Advanced Reporting

Welcome to Module 2\! Now that we have a working transaction system, we're going to add a critical feature for any real-world application: **background processing**.

When a user performs an action, they expect an immediate response. They shouldn't have to wait for the server to finish slow, non-essential tasks like generating reports. In this module, you'll implement a payment flow that immediately confirms the payment to the user, while a **goroutine** works in the background to update a complex reporting table.

### Learning Goals

  * **Understand Asynchronous Processing**: Learn why and when to run tasks in the background.
  * **Implement Goroutines**: Use the `go` keyword to spawn a new, concurrent process.
  * **Manage Graceful Shutdowns**: Use `sync.WaitGroup` to ensure your background jobs finish before the server shuts down.
  * **Prevent Race Conditions**: Use database locking to safely update shared data from multiple goroutines.
  * **Work with `JSON` in the Database**: Store and retrieve semi-structured data using a `JSON` column type and map it to a Go struct.

## Step 1: Database Migration - The Reporting Table

First, we need a dedicated table to hold our aggregated report data. This table will only ever have **one row**, which we will continuously update.

1.  **Create the Migration File**:
    Create a new file: `database/migrations/000006_create_transaction_reports_table.up.sql`

    ```sql
    CREATE TABLE transaction_reports (
        id TINYINT UNSIGNED PRIMARY KEY,
        total_revenue BIGINT UNSIGNED NOT NULL DEFAULT 0,
        total_paid_transactions BIGINT UNSIGNED NOT NULL DEFAULT 0,
        total_products_sold BIGINT UNSIGNED NOT NULL DEFAULT 0,
        total_unique_customers BIGINT UNSIGNED NOT NULL DEFAULT 0,
        category_summary JSON,
        updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
    );

    -- Insert the single row that we will always update.
    INSERT INTO transaction_reports (id) VALUES (1);
    ```

2.  **Create the `down` migration**:
    Create `database/migrations/000006_create_transaction_reports_table.down.sql`:

    ```sql
    DROP TABLE IF EXISTS transaction_reports;
    ```

3.  **Run the Migration**:

    ```bash
    docker-compose run --rm app go run ./cmd/migrate/main.go up
    ```

## Step 2: Model Creation - Mapping to Go

Now, let's create the Go `struct` for our new table. GORM has excellent support for the `JSON` type, which we can map directly to a Go `map` or `struct`.

Create a new file: `internal/model/transaction_report_model.go`

```go
package model

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
)

// CategorySummary defines the structure for our JSON data.
type CategorySummary map[string]int64

// Value makes CategorySummary implement the driver.Valuer interface.
// This is how GORM knows how to save our map into a JSON column.
func (cs CategorySummary) Value() (driver.Value, error) {
	return json.Marshal(cs)
}

// Scan makes CategorySummary implement the sql.Scanner interface.
// This is how GORM knows how to read data from a JSON column into our map.
func (cs *CategorySummary) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &cs)
}

// TransactionReport defines the model for our single-row reporting table.
type TransactionReport struct {
	ID                      uint8           `gorm:"primary_key"`
	TotalRevenue            uint64
	TotalPaidTransactions   uint64
	TotalProductsSold       uint64
	TotalUniqueCustomers    uint64
	CategorySummary         CategorySummary `gorm:"type:json"`
	UpdatedAt               time.Time
}

// Save updates the single report row.
func (tr *TransactionReport) Save(db *gorm.DB) error {
	return db.WithContext(context.Background()).Save(tr).Error
}
```

## Step 3: Service Layer - The Concurrent Logic

This is the core of the module. We will update our `TransactionService` to handle background processing.

1.  **Update `TransactionService`**:
    Open `internal/service/transaction_service.go` and modify the struct to include a `sync.WaitGroup`.

    ```go
    // Add sync to imports
    import "sync"

    type TransactionService struct {
    	db *gorm.DB
    	wg *sync.WaitGroup // <-- ADD THIS
    }

    // Update the constructor
    func NewTransactionService(db *gorm.DB, wg *sync.WaitGroup) *TransactionService {
    	return &TransactionService{db: db, wg: wg} // <-- ADD wg
    }
    ```

2.  **Update Service Initialization**:
    Open `internal/server/routes.go` and pass the `wg` when you create the service.

    ```go
    // In registerRoutes function
    transactionService := service.NewTransactionService(db, wg) // <-- PASS wg
    ```

3.  **Create the `MarkAsPaid` Method**:
    Add the following new method to `internal/service/transaction_service.go`.

    ```go
    // Add "gorm.io/gorm/clause" to imports

    // MarkAsPaid updates a transaction to paid and triggers a background report update.
    func (s *TransactionService) MarkAsPaid(ctx context.Context, transactionID uuid.UUID) error {
    	// --- SYNCHRONOUS PART ---
    	// The user waits for this to finish.
    	var transaction model.Transaction
    	if err := s.db.WithContext(ctx).First(&transaction, "id = ?", transactionID).Error; err != nil {
    		return errors.New("transaction not found")
    	}

    	isPaid := true
    	transaction.IsPaid = &isPaid
    	if err := transaction.Save(s.db); err != nil {
    		return err // Return error if the quick update fails
    	}

    	// --- ASYNCHRONOUS PART ---
    	// The user DOES NOT wait for this.
    	s.wg.Add(1) // Signal that a new background job has started.
    	go func() {
    		defer s.wg.Done() // Signal that the job is done when the function exits.

    		// Create a new background context.
    		bgCtx := context.Background()

    		// Use a database transaction to ensure all report updates succeed or fail together.
    		err := s.db.WithContext(bgCtx).Transaction(func(tx *gorm.DB) error {
    			// Lock the single report row to prevent race conditions from concurrent payments.
    			var report model.TransactionReport
    			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&report, "id = ?", 1).Error; err != nil {
    				return err
    			}

    			// --- Perform Aggregation Queries ---
    			var totalRevenue, totalPaidTransactions, totalUniqueCustomers uint64
    			tx.Model(&model.Transaction{}).Where("is_paid = ?", true).Select("COALESCE(SUM(total), 0)").Row().Scan(&totalRevenue)
    			tx.Model(&model.Transaction{}).Where("is_paid = ?", true).Count(&totalPaidTransactions)
    			tx.Model(&model.Transaction{}).Where("is_paid = ?", true).Select("COUNT(DISTINCT user_id)").Row().Scan(&totalUniqueCustomers)

    			var totalProductsSold uint64
    			tx.Model(&model.TransactionDetail{}).Joins("JOIN transactions ON transactions.id = transaction_details.transaction_id").
    				Where("transactions.is_paid = ?", true).Select("COALESCE(SUM(qty), 0)").Row().Scan(&totalProductsSold)

    			type CategoryResult struct {
    				Category model.ProductCategory
    				Count    int64
    			}
    			var categoryResults []CategoryResult
    			tx.Model(&model.TransactionDetail{}).Joins("JOIN transactions ON transactions.id = transaction_details.transaction_id").
    				Where("transactions.is_paid = ?", true).Select("category, SUM(qty) as count").Group("category").Find(&categoryResults)

    			// Update the report struct with the new values
    			report.TotalRevenue = totalRevenue
    			report.TotalPaidTransactions = totalPaidTransactions
    			report.TotalUniqueCustomers = totalUniqueCustomers
    			report.TotalProductsSold = totalProductsSold
    			report.CategorySummary = make(model.CategorySummary)
    			for _, res := range categoryResults {
    				// Convert category enum to string for JSON key
    				var catName string
    				switch res.Category {
    				case model.Goods: catName = "Goods"
    				case model.Service: catName = "Service"
    				case model.Subscription: catName = "Subscription"
    				}
    				report.CategorySummary[catName] = res.Count
    			}

    			// Save the updated report
    			return report.Save(tx)
    		})

    		if err != nil {
    			// In a real app, you would have a robust logging/alerting system here.
    			fmt.Printf("Error updating transaction report in background: %v\n", err)
    		}
    	}()

    	return nil
    }
    ```

## Step 4: Handler and Route

Finally, let's expose this functionality through an API endpoint.

1.  **Add Handler Method**:
    Open `internal/handler/http/transaction_handler.go` and add this new method.

    ```go
    // MarkAsPaid handles the request to mark a transaction as paid.
    // @Summary      Pay for a Transaction
    // @Description  Marks a transaction as paid and triggers a background report update.
    // @Tags         Transactions
    // @Produce      json
    // @Security     ApiKeyAuth
    // @Param        id   path      string  true  "Transaction ID"
    // @Success      200  {object}  response.ApiResponse "Successfully paid"
    // @Failure      401  {object}  response.ApiResponse "Unauthorized"
    // @Failure      404  {object}  response.ApiResponse "Transaction not found"
    // @Router       /transactions/{id}/pay [post]
    func (h *TransactionHandler) MarkAsPaid(c *fiber.Ctx) error {
    	idParam := c.Params("id")
    	transactionID, err := uuid.Parse(idParam)
    	if err != nil {
    		return response.Error(c, fiber.StatusBadRequest, errors.New("invalid ID format"))
    	}

    	err = h.transactionService.MarkAsPaid(c.Context(), transactionID)
    	if err != nil {
    		// Differentiate between not found and other errors
    		if strings.Contains(err.Error(), "not found") {
    			return response.Error(c, fiber.StatusNotFound, err)
    		}
    		return response.Error(c, fiber.StatusInternalServerError, err)
    	}

    	return response.Success(c, fiber.StatusOK, fiber.Map{"message": "Transaction marked as paid. Report is updating."})
    }
    ```

2.  **Register the Route**:
    Open `internal/server/routes.go` and add the new route.

    ```go
    // Inside the registerRoutes function
    // --- Transaction routes ---
    api.Post("/transactions", authMiddleware, transactionHandler.CreateTransaction)
    api.Post("/transactions/:id/pay", authMiddleware, transactionHandler.MarkAsPaid) // <-- ADD THIS
    ```

## Conclusion

You have now implemented a powerful, real-world feature. When the `/pay` endpoint is called, the user gets an instant response, and the heavy lifting of calculating and updating the report happens entirely in the background. You've learned how to use goroutines, manage them with a WaitGroup for safe shutdowns, and prevent data corruption with database locks.

-----