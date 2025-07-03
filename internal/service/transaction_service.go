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
