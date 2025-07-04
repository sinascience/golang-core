package service

import (
	"sync"
	"errors"
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"
	"venturo-core/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// TransactionService handles the business logic for transactions.
type TransactionService struct {
	db *gorm.DB
	wg *sync.WaitGroup
}

// NewTransactionService creates a new transaction service.
func NewTransactionService(db *gorm.DB, wg *sync.WaitGroup) *TransactionService {
	return &TransactionService{db: db, wg: wg}
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
			var totalRevenue, totalUniqueCustomers uint64
			var totalPaidTransactions int64
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
			report.TotalPaidTransactions = uint64(totalPaidTransactions)
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