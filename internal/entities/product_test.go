package entities

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNewProduct(t *testing.T) {
	catID := uuid.New()
	prod := NewProduct("Mouse", "Wireless mouse", 2999, "http://img.com", catID)

	if prod.Name != "Mouse" {
		t.Errorf("expected name %q, got %q", "Mouse", prod.Name)
	}
	if prod.Description != "Wireless mouse" {
		t.Errorf("expected description %q, got %q", "Wireless mouse", prod.Description)
	}
	if prod.Price != 2999 {
		t.Errorf("expected price 2999, got %d", prod.Price)
	}
	if prod.ImageURL != "http://img.com" {
		t.Errorf("expected image URL %q, got %q", "http://img.com", prod.ImageURL)
	}
	if prod.CategoryID != catID {
		t.Errorf("expected category ID %s, got %s", catID, prod.CategoryID)
	}

	if prod.ID.String() == "" {
		t.Error("expected non-empty ID")
	}
	if prod.CreatedAt.IsZero() {
		t.Error("expected non-zero CreatedAt")
	}
	if prod.UpdatedAt.IsZero() {
		t.Error("expected non-zero UpdatedAt")
	}
}

func TestProductResetUpdatedAt(t *testing.T) {
	prod := NewProduct("Keyboard", "Mechanical", 7999, "", uuid.Nil)
	old := prod.UpdatedAt

	time.Sleep(10 * time.Millisecond)
	prod.ResetUpdatedAt()

	if !prod.UpdatedAt.After(old) {
		t.Error("expected UpdatedAt to be updated")
	}
}
