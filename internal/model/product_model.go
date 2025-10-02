package model

import (
	"context"
	"time"
	"strconv"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)
type ProductType struct {
	Color string `json:"color"`
	Stock int16 `json:"stock"`
}

// Product defines the product model with image status.
type Product struct {
	ID           uuid.UUID `gorm:"type:char(36);primary_key"`
	Name         string    `gorm:"size:255;not null"`
	Slug         string	   `gorm:"size:255;not null;unique"`
	Price        int32
	Stock        int16
	ImageURL     string    `gorm:"size:255"`
	ImageStatus  string    `gorm:"size:20;not null;default:'default'"`
	Type         datatypes.JSON `gorm:"type:json"`
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
	var slug = p.Slug
	var count int

	for {
		// cek apakah slug sudah ada
		var existing Product
		result := db.WithContext(context.Background()).Where("slug = ?", p.Slug).First(&existing)

		if result.RowsAffected == 0 {
			// slug unik → langsung dipakai
			break
		}

		// slug sudah ada → tambahkan angka di belakang
		count++
		p.Slug = slug + "-" + strconv.Itoa(count)
	}

	return db.WithContext(context.Background()).Save(p).Error
}


// Get All Products
func (p *Product) GetAll(db *gorm.DB) ([]Product, error) {
	// return db.WithContext(context.Background()).First(p).Error
	var products []Product
	return products, db.WithContext(context.Background()).Find(&products).Error
}

// FindByID retrieves a single post by
func (p *Product) FindByID(db *gorm.DB, id uuid.UUID) (*Product, error) {
	var product Product
	err := db.WithContext(context.Background()).Where("id = ?", id).First(&product).Error
	return &product, err
}

// FindBySlug retrieves a single post by
func (p *Product) FindBySlug(db *gorm.DB, slug string) (*Product, error) {
	var product Product
	err := db.WithContext(context.Background()).Where("slug = ?", slug).First(&product).Error
	return &product, err
}