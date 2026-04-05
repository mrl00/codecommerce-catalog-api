package entities

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	ImageURL    string    `json:"image_url"`
	CategoryID  uuid.UUID `json:"category"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func NewProduct(name string, description string, price float64, imageURL string, categoryID uuid.UUID) *Product {
	now := time.Now()
	return &Product{
		ID:          uuid.New(),
		Name:        name,
		Description: description,
		Price:       price,
		ImageURL:    imageURL,
		CategoryID:  categoryID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}
