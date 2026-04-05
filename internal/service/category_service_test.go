package service

import (
	"testing"

	"codecommerceapi/internal/entities"

	"github.com/google/uuid"
)

type mockCategoryRepo struct {
	saved      []*entities.Category
	byID       map[string]*entities.Category
	all        []*entities.Category
	totalCount int
	deleteErr  error
}

func (m *mockCategoryRepo) SaveCategory(c *entities.Category) error {
	m.saved = append(m.saved, c)
	return nil
}

func (m *mockCategoryRepo) FindCategoryByID(id uuid.UUID) (*entities.Category, error) {
	c, ok := m.byID[id.String()]
	if !ok {
		return nil, ErrCategoryNotFound
	}
	return c, nil
}

func (m *mockCategoryRepo) FindAllCategories(params PaginationParams) (*PaginatedResult[*entities.Category], error) {
	total := len(m.all)
	start := (params.Page - 1) * params.PerPage
	end := start + params.PerPage
	if end > total {
		end = total
	}

	var items []*entities.Category
	if start < total {
		items = m.all[start:end]
	} else {
		items = []*entities.Category{}
	}

	totalPages := (total + params.PerPage - 1) / params.PerPage
	return &PaginatedResult[*entities.Category]{
		Items:      items,
		Total:      total,
		Page:       params.Page,
		PerPage:    params.PerPage,
		TotalPages: totalPages,
	}, nil
}

func (m *mockCategoryRepo) UpdateCategory(c *entities.Category) error {
	m.byID[c.ID.String()] = c
	return nil
}

func (m *mockCategoryRepo) DeleteCategory(id uuid.UUID) error {
	delete(m.byID, id.String())
	return m.deleteErr
}

func newTestCategoryService(repo *mockCategoryRepo) *CategoryService {
	if repo == nil {
		repo = &mockCategoryRepo{byID: make(map[string]*entities.Category)}
	}
	return &CategoryService{repo: repo}
}

func TestCategoryService_CreateCategory_Success(t *testing.T) {
	repo := &mockCategoryRepo{byID: make(map[string]*entities.Category)}
	svc := newTestCategoryService(repo)

	cat, err := svc.CreateCategory("  Electronics  ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cat.Name != "Electronics" {
		t.Errorf("expected trimmed name, got %q", cat.Name)
	}

	if len(repo.saved) != 1 {
		t.Error("expected category to be saved")
	}
}

func TestCategoryService_CreateCategory_EmptyName(t *testing.T) {
	svc := newTestCategoryService(nil)

	_, err := svc.CreateCategory("   ")
	if err == nil {
		t.Fatal("expected error for empty name")
	}

	if err.Error() != "category name cannot be empty" {
		t.Errorf("got wrong error: %v", err)
	}
}

func TestCategoryService_GetCategory_Success(t *testing.T) {
	existing := entities.NewCategory("Books")
	repo := &mockCategoryRepo{byID: map[string]*entities.Category{existing.ID.String(): existing}}
	svc := newTestCategoryService(repo)

	cat, err := svc.GetCategory(existing.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cat.Name != "Books" {
		t.Errorf("expected name Books, got %q", cat.Name)
	}
}

func TestCategoryService_GetCategory_NotFound(t *testing.T) {
	repo := &mockCategoryRepo{byID: make(map[string]*entities.Category)}
	svc := newTestCategoryService(repo)

	_, err := svc.GetCategory(uuid.New())
	if err != ErrCategoryNotFound {
		t.Errorf("expected ErrCategoryNotFound, got %v", err)
	}
}

func TestCategoryService_UpdateCategory_Success(t *testing.T) {
	existing := entities.NewCategory("Old")
	repo := &mockCategoryRepo{byID: map[string]*entities.Category{existing.ID.String(): existing}}
	svc := newTestCategoryService(repo)

	updated, err := svc.UpdateCategory(existing.ID, "New Name")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if updated.Name != "New Name" {
		t.Errorf("expected name New Name, got %q", updated.Name)
	}
}

func TestCategoryService_UpdateCategory_NotFound(t *testing.T) {
	svc := newTestCategoryService(&mockCategoryRepo{byID: make(map[string]*entities.Category)})

	_, err := svc.UpdateCategory(uuid.New(), "Name")
	if err != ErrCategoryNotFound {
		t.Errorf("expected ErrCategoryNotFound, got %v", err)
	}
}

func TestCategoryService_DeleteCategory_Success(t *testing.T) {
	existing := entities.NewCategory("To Delete")
	repo := &mockCategoryRepo{
		byID: map[string]*entities.Category{existing.ID.String(): existing},
	}
	svc := newTestCategoryService(repo)

	err := svc.DeleteCategory(existing.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, ok := repo.byID[existing.ID.String()]; ok {
		t.Error("expected category to be deleted")
	}
}

func TestCategoryService_ListCategories_Pagination(t *testing.T) {
	all := []*entities.Category{
		entities.NewCategory("A"),
		entities.NewCategory("B"),
		entities.NewCategory("C"),
	}
	repo := &mockCategoryRepo{all: all, byID: make(map[string]*entities.Category)}
	svc := newTestCategoryService(repo)

	result, err := svc.ListCategories(PaginationParams{Page: 1, PerPage: 2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Items) != 2 {
		t.Errorf("expected 2 items, got %d", len(result.Items))
	}
	if result.Total != 3 {
		t.Errorf("expected total 3, got %d", result.Total)
	}
}
