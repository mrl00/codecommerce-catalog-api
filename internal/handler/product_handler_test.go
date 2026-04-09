package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"codecommerceapi/internal/entities"
	"codecommerceapi/internal/service"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// ---------------------------------------------------------------------------
// mock ProductRepository
// ---------------------------------------------------------------------------

type mockProductRepo struct {
	products  map[string]*entities.Product
	saveErr   error
	deleteErr error
}

func newMockProductRepo() *mockProductRepo {
	return &mockProductRepo{products: make(map[string]*entities.Product)}
}

func (m *mockProductRepo) SaveProduct(p *entities.Product) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	m.products[p.ID.String()] = p
	return nil
}

func (m *mockProductRepo) FindProductByID(id uuid.UUID) (*entities.Product, error) {
	p, ok := m.products[id.String()]
	if !ok {
		return nil, service.ErrProductNotFound
	}
	return p, nil
}

func (m *mockProductRepo) FindAllProducts(params service.PaginationParams) (*service.PaginatedResult[*entities.Product], error) {
	all := make([]*entities.Product, 0, len(m.products))
	for _, p := range m.products {
		all = append(all, p)
	}
	total := len(all)
	start := (params.Page - 1) * params.PerPage
	if start > total {
		start = total
	}
	end := start + params.PerPage
	if end > total {
		end = total
	}
	totalPages := 0
	if total > 0 {
		totalPages = (total + params.PerPage - 1) / params.PerPage
	}
	return &service.PaginatedResult[*entities.Product]{
		Items:      all[start:end],
		Total:      total,
		Page:       params.Page,
		PerPage:    params.PerPage,
		TotalPages: totalPages,
	}, nil
}

func (m *mockProductRepo) FindProductsByCategoryID(categoryID uuid.UUID, params service.PaginationParams) (*service.PaginatedResult[*entities.Product], error) {
	filtered := make([]*entities.Product, 0)
	for _, p := range m.products {
		if p.CategoryID == categoryID {
			filtered = append(filtered, p)
		}
	}
	total := len(filtered)
	start := (params.Page - 1) * params.PerPage
	if start > total {
		start = total
	}
	end := start + params.PerPage
	if end > total {
		end = total
	}
	totalPages := 0
	if total > 0 {
		totalPages = (total + params.PerPage - 1) / params.PerPage
	}
	return &service.PaginatedResult[*entities.Product]{
		Items:      filtered[start:end],
		Total:      total,
		Page:       params.Page,
		PerPage:    params.PerPage,
		TotalPages: totalPages,
	}, nil
}

func (m *mockProductRepo) UpdateProduct(p *entities.Product) error {
	m.products[p.ID.String()] = p
	return nil
}

func (m *mockProductRepo) DeleteProduct(id uuid.UUID) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	if _, ok := m.products[id.String()]; !ok {
		return service.ErrProductNotFound
	}
	delete(m.products, id.String())
	return nil
}

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

func newTestProductHandler(repo *mockProductRepo) *ProductHandler {
	return NewProductHandler(service.NewProductService(repo))
}

func productRouter(h *ProductHandler) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/api/products", h.CreateProduct).Methods("POST")
	r.HandleFunc("/api/products", h.ListProducts).Methods("GET")
	r.HandleFunc("/api/products/{id}", h.GetProduct).Methods("GET")
	r.HandleFunc("/api/products/{id}", h.UpdateProduct).Methods("PUT")
	r.HandleFunc("/api/products/{id}", h.DeleteProduct).Methods("DELETE")
	r.HandleFunc("/api/categories/{id}/products", h.ListProductsByCategory).Methods("GET")
	return r
}

func seedProduct(repo *mockProductRepo, name string, price int64, categoryID uuid.UUID) *entities.Product {
	p := entities.NewProduct(name, "test description", price, "http://img.test/photo.jpg", categoryID)
	repo.products[p.ID.String()] = p
	return p
}

func validProductJSON(categoryID uuid.UUID) string {
	return fmt.Sprintf(`{
		"name":"Wireless Mouse",
		"description":"Ergonomic wireless mouse",
		"price":2999,
		"image_url":"http://img.test/mouse.jpg",
		"category":"%s"
	}`, categoryID)
}

// ---------------------------------------------------------------------------
// parsePaginationParams
// ---------------------------------------------------------------------------

func TestParsePaginationParams(t *testing.T) {
	t.Run("defaults when no query params", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/products", nil)
		p := parsePaginationParams(req)
		if p.Page != 1 {
			t.Errorf("expected page 1, got %d", p.Page)
		}
		if p.PerPage != 10 {
			t.Errorf("expected per_page 10, got %d", p.PerPage)
		}
	})

	t.Run("custom values", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/products?page=5&per_page=50", nil)
		p := parsePaginationParams(req)
		if p.Page != 5 {
			t.Errorf("expected page 5, got %d", p.Page)
		}
		if p.PerPage != 50 {
			t.Errorf("expected per_page 50, got %d", p.PerPage)
		}
	})

	t.Run("page below 1 defaults to 1", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/products?page=-1", nil)
		p := parsePaginationParams(req)
		if p.Page != 1 {
			t.Errorf("expected page 1, got %d", p.Page)
		}
	})

	t.Run("per_page above 100 defaults to 10", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/products?per_page=101", nil)
		p := parsePaginationParams(req)
		if p.PerPage != 10 {
			t.Errorf("expected per_page 10, got %d", p.PerPage)
		}
	})

	t.Run("per_page below 1 defaults to 10", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/products?per_page=0", nil)
		p := parsePaginationParams(req)
		if p.PerPage != 10 {
			t.Errorf("expected per_page 10, got %d", p.PerPage)
		}
	})

	t.Run("non-numeric values use defaults", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/products?page=abc&per_page=xyz", nil)
		p := parsePaginationParams(req)
		if p.Page != 1 {
			t.Errorf("expected page 1, got %d", p.Page)
		}
		if p.PerPage != 10 {
			t.Errorf("expected per_page 10, got %d", p.PerPage)
		}
	})
}

// ---------------------------------------------------------------------------
// CreateProduct
// ---------------------------------------------------------------------------

func TestCreateProduct(t *testing.T) {
	catID := uuid.New()

	t.Run("success", func(t *testing.T) {
		repo := newMockProductRepo()
		r := productRouter(newTestProductHandler(repo))

		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/products",
			bytes.NewBufferString(validProductJSON(catID)))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(rr, req)

		assertStatus(t, rr, http.StatusCreated)
		assertContentType(t, rr, "application/json")

		var resp entities.Product
		if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if resp.Name != "Wireless Mouse" {
			t.Errorf("expected name %q, got %q", "Wireless Mouse", resp.Name)
		}
		if resp.Price != 2999 {
			t.Errorf("expected price 2999, got %d", resp.Price)
		}
		if resp.ID == uuid.Nil {
			t.Error("expected non-nil ID")
		}
		if len(repo.products) != 1 {
			t.Errorf("expected 1 saved product, got %d", len(repo.products))
		}
	})

	t.Run("invalid JSON body", func(t *testing.T) {
		repo := newMockProductRepo()
		r := productRouter(newTestProductHandler(repo))

		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/products",
			bytes.NewBufferString(`{invalid}`))
		r.ServeHTTP(rr, req)

		assertStatus(t, rr, http.StatusBadRequest)
		assertJSONError(t, rr, "invalid request body")
	})

	t.Run("empty name", func(t *testing.T) {
		repo := newMockProductRepo()
		r := productRouter(newTestProductHandler(repo))

		rr := httptest.NewRecorder()
		body := fmt.Sprintf(`{"name":"","price":100,"category":"%s"}`, catID)
		req := httptest.NewRequest(http.MethodPost, "/api/products",
			bytes.NewBufferString(body))
		r.ServeHTTP(rr, req)

		assertStatus(t, rr, http.StatusBadRequest)
		assertJSONError(t, rr, "name is required")
	})

	t.Run("zero price", func(t *testing.T) {
		repo := newMockProductRepo()
		r := productRouter(newTestProductHandler(repo))

		rr := httptest.NewRecorder()
		body := fmt.Sprintf(`{"name":"Item","price":0,"category":"%s"}`, catID)
		req := httptest.NewRequest(http.MethodPost, "/api/products",
			bytes.NewBufferString(body))
		r.ServeHTTP(rr, req)

		assertStatus(t, rr, http.StatusBadRequest)
		assertJSONError(t, rr, "price must be greater than zero")
	})

	t.Run("negative price", func(t *testing.T) {
		repo := newMockProductRepo()
		r := productRouter(newTestProductHandler(repo))

		rr := httptest.NewRecorder()
		body := fmt.Sprintf(`{"name":"Item","price":-10,"category":"%s"}`, catID)
		req := httptest.NewRequest(http.MethodPost, "/api/products",
			bytes.NewBufferString(body))
		r.ServeHTTP(rr, req)

		assertStatus(t, rr, http.StatusBadRequest)
		assertJSONError(t, rr, "price must be greater than zero")
	})

	t.Run("missing category", func(t *testing.T) {
		repo := newMockProductRepo()
		r := productRouter(newTestProductHandler(repo))

		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/products",
			bytes.NewBufferString(`{"name":"Item","price":100}`))
		r.ServeHTTP(rr, req)

		assertStatus(t, rr, http.StatusBadRequest)
		assertJSONError(t, rr, "category is required")
	})

	t.Run("invalid category ID format", func(t *testing.T) {
		repo := newMockProductRepo()
		r := productRouter(newTestProductHandler(repo))

		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/products",
			bytes.NewBufferString(`{"name":"Item","price":100,"category":"not-a-uuid"}`))
		r.ServeHTTP(rr, req)

		assertStatus(t, rr, http.StatusBadRequest)
		assertJSONError(t, rr, "invalid category id format")
	})
}

// ---------------------------------------------------------------------------
// GetProduct
// ---------------------------------------------------------------------------

func TestGetProduct(t *testing.T) {
	catID := uuid.New()

	t.Run("success", func(t *testing.T) {
		repo := newMockProductRepo()
		prod := seedProduct(repo, "Mouse", 2999, catID)
		r := productRouter(newTestProductHandler(repo))

		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/products/%s", prod.ID), nil)
		r.ServeHTTP(rr, req)

		assertStatus(t, rr, http.StatusOK)
		assertContentType(t, rr, "application/json")

		var resp entities.Product
		if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if resp.Name != "Mouse" {
			t.Errorf("expected name %q, got %q", "Mouse", resp.Name)
		}
		if resp.Price != 2999 {
			t.Errorf("expected price 2999, got %d", resp.Price)
		}
	})

	t.Run("invalid ID", func(t *testing.T) {
		repo := newMockProductRepo()
		r := productRouter(newTestProductHandler(repo))

		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/products/not-a-uuid", nil)
		r.ServeHTTP(rr, req)

		assertStatus(t, rr, http.StatusBadRequest)
		assertJSONError(t, rr, "invalid product id")
	})

	t.Run("not found", func(t *testing.T) {
		repo := newMockProductRepo()
		r := productRouter(newTestProductHandler(repo))

		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/products/%s", uuid.New()), nil)
		r.ServeHTTP(rr, req)

		assertStatus(t, rr, http.StatusNotFound)
		assertJSONError(t, rr, "product not found")
	})
}

// ---------------------------------------------------------------------------
// ListProducts
// ---------------------------------------------------------------------------

func TestListProducts(t *testing.T) {
	catID := uuid.New()

	t.Run("returns paginated result", func(t *testing.T) {
		repo := newMockProductRepo()
		seedProduct(repo, "Mouse", 2999, catID)
		seedProduct(repo, "Keyboard", 7999, catID)
		r := productRouter(newTestProductHandler(repo))

		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/products?page=1&per_page=10", nil)
		r.ServeHTTP(rr, req)

		assertStatus(t, rr, http.StatusOK)
		assertContentType(t, rr, "application/json")

		var resp service.PaginatedResult[json.RawMessage]
		if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if resp.Total != 2 {
			t.Errorf("expected total 2, got %d", resp.Total)
		}
		if resp.Page != 1 {
			t.Errorf("expected page 1, got %d", resp.Page)
		}
	})

	t.Run("empty list", func(t *testing.T) {
		repo := newMockProductRepo()
		r := productRouter(newTestProductHandler(repo))

		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/products", nil)
		r.ServeHTTP(rr, req)

		assertStatus(t, rr, http.StatusOK)

		var resp service.PaginatedResult[json.RawMessage]
		if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if resp.Total != 0 {
			t.Errorf("expected total 0, got %d", resp.Total)
		}
	})
}

// ---------------------------------------------------------------------------
// ListProductsByCategory
// ---------------------------------------------------------------------------

func TestListProductsByCategory(t *testing.T) {
	catA := uuid.New()
	catB := uuid.New()

	t.Run("filters by category", func(t *testing.T) {
		repo := newMockProductRepo()
		seedProduct(repo, "Mouse", 2999, catA)
		seedProduct(repo, "Keyboard", 7999, catA)
		seedProduct(repo, "Novel", 1599, catB)
		r := productRouter(newTestProductHandler(repo))

		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet,
			fmt.Sprintf("/api/categories/%s/products?page=1&per_page=10", catA), nil)
		r.ServeHTTP(rr, req)

		assertStatus(t, rr, http.StatusOK)
		assertContentType(t, rr, "application/json")

		var resp service.PaginatedResult[json.RawMessage]
		if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if resp.Total != 2 {
			t.Errorf("expected total 2 for catA, got %d", resp.Total)
		}
	})

	t.Run("invalid category ID", func(t *testing.T) {
		repo := newMockProductRepo()
		r := productRouter(newTestProductHandler(repo))

		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/categories/bad-id/products", nil)
		r.ServeHTTP(rr, req)

		assertStatus(t, rr, http.StatusBadRequest)
		assertJSONError(t, rr, "invalid category id")
	})

	t.Run("empty result for unknown category", func(t *testing.T) {
		repo := newMockProductRepo()
		r := productRouter(newTestProductHandler(repo))

		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet,
			fmt.Sprintf("/api/categories/%s/products", uuid.New()), nil)
		r.ServeHTTP(rr, req)

		assertStatus(t, rr, http.StatusOK)

		var resp service.PaginatedResult[json.RawMessage]
		if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if resp.Total != 0 {
			t.Errorf("expected total 0, got %d", resp.Total)
		}
	})
}

// ---------------------------------------------------------------------------
// UpdateProduct
// ---------------------------------------------------------------------------

func TestUpdateProduct(t *testing.T) {
	catID := uuid.New()

	t.Run("success", func(t *testing.T) {
		repo := newMockProductRepo()
		prod := seedProduct(repo, "Old Mouse", 1999, catID)
		r := productRouter(newTestProductHandler(repo))

		body := fmt.Sprintf(`{
			"name":"New Mouse",
			"description":"Updated",
			"price":3499,
			"image_url":"http://img.test/new.jpg",
			"category":"%s"
		}`, catID)
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut,
			fmt.Sprintf("/api/products/%s", prod.ID),
			bytes.NewBufferString(body))
		r.ServeHTTP(rr, req)

		assertStatus(t, rr, http.StatusOK)
		assertContentType(t, rr, "application/json")

		var resp entities.Product
		if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if resp.Name != "New Mouse" {
			t.Errorf("expected name %q, got %q", "New Mouse", resp.Name)
		}
		if resp.Price != 3499 {
			t.Errorf("expected price 3499, got %d", resp.Price)
		}
	})

	t.Run("invalid ID", func(t *testing.T) {
		repo := newMockProductRepo()
		r := productRouter(newTestProductHandler(repo))

		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, "/api/products/bad-id",
			bytes.NewBufferString(validProductJSON(catID)))
		r.ServeHTTP(rr, req)

		assertStatus(t, rr, http.StatusBadRequest)
		assertJSONError(t, rr, "invalid product id")
	})

	t.Run("invalid JSON body", func(t *testing.T) {
		repo := newMockProductRepo()
		prod := seedProduct(repo, "Mouse", 2999, catID)
		r := productRouter(newTestProductHandler(repo))

		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut,
			fmt.Sprintf("/api/products/%s", prod.ID),
			bytes.NewBufferString(`{bad}`))
		r.ServeHTTP(rr, req)

		assertStatus(t, rr, http.StatusBadRequest)
		assertJSONError(t, rr, "invalid request body")
	})

	t.Run("empty name", func(t *testing.T) {
		repo := newMockProductRepo()
		prod := seedProduct(repo, "Mouse", 2999, catID)
		r := productRouter(newTestProductHandler(repo))

		body := fmt.Sprintf(`{"name":"","price":100,"category":"%s"}`, catID)
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut,
			fmt.Sprintf("/api/products/%s", prod.ID),
			bytes.NewBufferString(body))
		r.ServeHTTP(rr, req)

		assertStatus(t, rr, http.StatusBadRequest)
		assertJSONError(t, rr, "name is required")
	})

	t.Run("zero price", func(t *testing.T) {
		repo := newMockProductRepo()
		prod := seedProduct(repo, "Mouse", 2999, catID)
		r := productRouter(newTestProductHandler(repo))

		body := fmt.Sprintf(`{"name":"Mouse","price":0,"category":"%s"}`, catID)
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut,
			fmt.Sprintf("/api/products/%s", prod.ID),
			bytes.NewBufferString(body))
		r.ServeHTTP(rr, req)

		assertStatus(t, rr, http.StatusBadRequest)
		assertJSONError(t, rr, "price must be greater than zero")
	})

	t.Run("missing category", func(t *testing.T) {
		repo := newMockProductRepo()
		prod := seedProduct(repo, "Mouse", 2999, catID)
		r := productRouter(newTestProductHandler(repo))

		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut,
			fmt.Sprintf("/api/products/%s", prod.ID),
			bytes.NewBufferString(`{"name":"Mouse","price":100}`))
		r.ServeHTTP(rr, req)

		assertStatus(t, rr, http.StatusBadRequest)
		assertJSONError(t, rr, "category is required")
	})

	t.Run("invalid category ID format", func(t *testing.T) {
		repo := newMockProductRepo()
		prod := seedProduct(repo, "Mouse", 2999, catID)
		r := productRouter(newTestProductHandler(repo))

		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut,
			fmt.Sprintf("/api/products/%s", prod.ID),
			bytes.NewBufferString(`{"name":"Mouse","price":100,"category":"not-a-uuid"}`))
		r.ServeHTTP(rr, req)

		assertStatus(t, rr, http.StatusBadRequest)
		assertJSONError(t, rr, "invalid category id format")
	})

	t.Run("not found", func(t *testing.T) {
		repo := newMockProductRepo()
		r := productRouter(newTestProductHandler(repo))

		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut,
			fmt.Sprintf("/api/products/%s", uuid.New()),
			bytes.NewBufferString(validProductJSON(catID)))
		r.ServeHTTP(rr, req)

		assertStatus(t, rr, http.StatusNotFound)
		assertJSONError(t, rr, "product not found")
	})
}

// ---------------------------------------------------------------------------
// DeleteProduct
// ---------------------------------------------------------------------------

func TestDeleteProduct(t *testing.T) {
	catID := uuid.New()

	t.Run("success", func(t *testing.T) {
		repo := newMockProductRepo()
		prod := seedProduct(repo, "To Delete", 999, catID)
		r := productRouter(newTestProductHandler(repo))

		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/products/%s", prod.ID), nil)
		r.ServeHTTP(rr, req)

		assertStatus(t, rr, http.StatusNoContent)

		if _, ok := repo.products[prod.ID.String()]; ok {
			t.Error("expected product to be deleted from repo")
		}
	})

	t.Run("invalid ID", func(t *testing.T) {
		repo := newMockProductRepo()
		r := productRouter(newTestProductHandler(repo))

		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/api/products/bad-id", nil)
		r.ServeHTTP(rr, req)

		assertStatus(t, rr, http.StatusBadRequest)
		assertJSONError(t, rr, "invalid product id")
	})

	t.Run("not found", func(t *testing.T) {
		repo := newMockProductRepo()
		r := productRouter(newTestProductHandler(repo))

		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/products/%s", uuid.New()), nil)
		r.ServeHTTP(rr, req)

		assertStatus(t, rr, http.StatusNotFound)
		assertJSONError(t, rr, "product not found")
	})
}
