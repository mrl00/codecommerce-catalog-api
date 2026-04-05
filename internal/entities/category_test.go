package entities

import (
	"testing"
	"time"
)

func TestNewCategory(t *testing.T) {
	name := "Electronics"
	cat := NewCategory(name)

	if cat.Name != name {
		t.Errorf("expected name %q, got %q", name, cat.Name)
	}

	if cat.ID.String() == "" {
		t.Error("expected non-empty ID")
	}

	if cat.CreatedAt.IsZero() {
		t.Error("expected non-zero CreatedAt")
	}

	if cat.UpdatedAt.IsZero() {
		t.Error("expected non-zero UpdatedAt")
	}
}

func TestCategoryResetUpdatedAt(t *testing.T) {
	cat := NewCategory("Books")
	old := cat.UpdatedAt

	time.Sleep(10 * time.Millisecond)
	cat.ResetUpdatedAt()

	if !cat.UpdatedAt.After(old) {
		t.Error("expected UpdatedAt to be updated")
	}
}
