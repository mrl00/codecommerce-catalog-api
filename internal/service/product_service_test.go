package service

import (
	"testing"

	"codecommerceapi/internal/entities"

	"github.com/google/uuid"
)

type mockProductRepo struct {
	saved      []*entities.Product
	byID       map[string]*entities.Product
	all        []*entities.Product
	deleteErr  error
}

func (m *mockProductRepo) SaveProduct(p *entities.Product) error {
	m.saved = append(m.saved, p)
	return nil
}

func (m *mockProductRepo) FindProductByID(id uuid.UUID) (*entities.Product, error) {
	p, ok := m.byID[id.String()]
	if !ok {
		return nil, ErrProductNotFound
	}
	return p, nil
}

func (m *mockProductRepo) FindAllProducts(params PaginationParams) (*PaginatedResult[*entities.Product], error) {
	total := len(m.all)
	start := (params.Page - 1) * params.PerPage
	end := start + params.PerPage
	if end > total {
		end = total
	}

	var items []*entities.Product
	if start < total {
		items = m.all[start:end]
	} else {
		items = []*entities.Product{}
	}

	totalPages := (total + params.PerPage - 1) / params.PerPage
	return &PaginatedResult[*entities.Product]{
		Items:      items,
		Total:      total,
		Page:       params.Page,
		PerPage:    params.PerPage,
		TotalPages: totalPages,
	}, nil
}

func (m *mockProductRepo) FindProductsByCategoryID(categoryID uuid.UUID, params PaginationParams) (*PaginatedResult[*entities.Product], error) {
	var filtered []*entities.Product
	for _, p := range m.all {
		if p.CategoryID == categoryID {
			filtered = append(filtered, p)
		}
	}

	total := len(filtered)
	start := (params.Page - 1) * params.PerPage
	end := start + params.PerPage
	if end > total {
		end = total
	}

	var items []*entities.Product
	if start < total {
		items = filtered[start:end]
	} else {
		items = []*entities.Product{}
	}

	totalPages := (total + params.PerPage - 1) / params.PerPage
	return &PaginatedResult[*entities.Product]{
		Items:      items,
		Total:      total,
		Page:       params.Page,
		PerPage:    params.PerPage,
		TotalPages: totalPages,
	}, nil
}

func (m *mockProductRepo) UpdateProduct(p *entities.Product) error {
	m.byID[p.ID.String()] = p
	return nil
}

func (m *mockProductRepo) DeleteProduct(id uuid.UUID) error {
	delete(m.byID, id.String())
	return m.deleteErr
}

func newTestProductService(repo *mockProductRepo) *ProductService {
	if repo == nil {
		repo = &mockProductRepo{byID: make(map[string]*entities.Product)}
	}
	return &ProductService{repo: repo}
}

func TestProductService_CreateProduct_Success(t *testing.T) {
	repo := &mockProductRepo{byID: make(map[string]*entities.Product)}
	svc := newTestProductService(repo)

	catID := uuid.New()
	prod, err := svc.CreateProduct("Mouse", "Wireless mouse", 2999, "http://img.com", catID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if prod.Name != "Mouse" {
		t.Errorf("expected name Mouse, got %q", prod.Name)
	}
	if prod.Price != 2999 {
		t.Errorf("expected price 2999, got %d", prod.Price)
	}
	if len(repo.saved) != 1 {
		t.Error("expected product to be saved")
	}
}

func TestProductService_CreateProduct_EmptyName(t *testing.T) {
	svc := newTestProductService(nil)

	_, err := svc.CreateProduct("", "desc", 100, "", uuid.New())
	if err == nil {
		t.Fatal("expected error for empty name")
	}
	if err.Error() != "product name cannot be empty" {
		t.Errorf("got wrong error: %v", err)
	}
}

func TestProductService_CreateProduct_NegativePrice(t *testing.T) {
	svc := newTestProductService(nil)

	_, err := svc.CreateProduct("Product", "desc", -10, "", uuid.New())
	if err == nil {
		t.Fatal("expected error for negative price")
	}
	if err.Error() != "product price cannot be negative" {
		t.Errorf("got wrong error: %v", err)
	}
}

func TestProductService_CreateProduct_NilCategoryID(t *testing.T) {
	svc := newTestProductService(nil)

	_, err := svc.CreateProduct("Product", "desc", 100, "", uuid.Nil)
	if err == nil {
		t.Fatal("expected error for nil category ID")
	}
	if err.Error() != "product category ID is required" {
		t.Errorf("got wrong error: %v", err)
	}
}

func TestProductService_GetProduct_Success(t *testing.T) {
	catID := uuid.New()
	existing := entities.NewProduct("Mouse", "Wireless", 2999, "", catID)
	repo := &mockProductRepo{byID: map[string]*entities.Product{existing.ID.String(): existing}}
	svc := newTestProductService(repo)

	prod, err := svc.GetProduct(existing.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if prod.Name != "Mouse" {
		t.Errorf("expected name Mouse, got %q", prod.Name)
	}
}

func TestProductService_GetProduct_NotFound(t *testing.T) {
	repo := &mockProductRepo{byID: make(map[string]*entities.Product)}
	svc := newTestProductService(repo)

	_, err := svc.GetProduct(uuid.New())
	if err != ErrProductNotFound {
		t.Errorf("expected ErrProductNotFound, got %v", err)
	}
}

func TestProductService_UpdateProduct_Success(t *testing.T) {
	catID := uuid.New()
	existing := entities.NewProduct("Old", "Old desc", 1000, "", catID)
	repo := &mockProductRepo{byID: map[string]*entities.Product{existing.ID.String(): existing}}
	svc := newTestProductService(repo)

	updated, err := svc.UpdateProduct(existing.ID, "New Name", "New desc", 2000, "http://img", catID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if updated.Name != "New Name" {
		t.Errorf("expected name New Name, got %q", updated.Name)
	}
	if updated.Price != 2000 {
		t.Errorf("expected price 2000, got %d", updated.Price)
	}
}

func TestProductService_UpdateProduct_NotFound(t *testing.T) {
	repo := &mockProductRepo{byID: make(map[string]*entities.Product)}
	svc := newTestProductService(repo)

	_, err := svc.UpdateProduct(uuid.New(), "Name", "desc", 100, "", uuid.New())
	if err != ErrProductNotFound {
		t.Errorf("expected ErrProductNotFound, got %v", err)
	}
}

func TestProductService_DeleteProduct_Success(t *testing.T) {
	catID := uuid.New()
	existing := entities.NewProduct("Delete", "desc", 100, "", catID)
	repo := &mockProductRepo{
		byID: map[string]*entities.Product{existing.ID.String(): existing},
	}
	svc := newTestProductService(repo)

	err := svc.DeleteProduct(existing.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, ok := repo.byID[existing.ID.String()]; ok {
		t.Error("expected product to be deleted")
	}
}

func TestProductService_ListProducts_Pagination(t *testing.T) {
	catID := uuid.New()
	all := []*entities.Product{
		entities.NewProduct("A", "", 100, "", catID),
		entities.NewProduct("B", "", 200, "", catID),
		entities.NewProduct("C", "", 300, "", catID),
	}
	repo := &mockProductRepo{all: all, byID: make(map[string]*entities.Product)}
	svc := newTestProductService(repo)

	result, err := svc.ListProducts(PaginationParams{Page: 1, PerPage: 2})
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

func TestProductService_ListProductsByCategory(t *testing.T) {
	catA := uuid.New()
	catB := uuid.New()
	all := []*entities.Product{
		entities.NewProduct("A1", "", 100, "", catA),
		entities.NewProduct("A2", "", 200, "", catA),
		entities.NewProduct("B1", "", 300, "", catB),
	}
	repo := &mockProductRepo{all: all, byID: make(map[string]*entities.Product)}
	svc := newTestProductService(repo)

	result, err := svc.ListProductsByCategory(catA, PaginationParams{Page: 1, PerPage: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Total != 2 {
		t.Errorf("expected total 2, got %d", result.Total)
	}
}
