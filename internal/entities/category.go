package entities

import (
	"time"

	"github.com/google/uuid"
)

type Category struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewCategory(name string) *Category {
	now := time.Now()
	return &Category{
		ID:        uuid.New(),
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (c *Category) ResetUpdatedAt() {
	c.UpdatedAt = time.Now()
}
