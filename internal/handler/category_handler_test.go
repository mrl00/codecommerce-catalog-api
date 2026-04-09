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
// mock CategoryRepository
// ---------------------------------------------------------------------------

type mockCategoryRepo struct {
	categories map[string]*entities.Category
	saveErr    error
	deleteErr  error
}

func newMockCategoryRepo() *mockCategoryRepo {
	return &mockCategoryRepo{categories: make(map[string]*entities.Category)}
}

func (m *mockCategoryRepo) SaveCategory(c *entities.Category) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	m.categories[c.ID.String()] = c
	return nil
}

func (m *mockCategoryRepo) FindCategoryByID(id uuid.UUID) (*entities.Category, error) {
	c, ok := m.categories[id.String()]
	if !ok {
		return nil, service.ErrCategoryNotFound
	}
	return c, nil
}

func (m *mockCategoryRepo) FindAllCategories(params service.PaginationParams) (*service.PaginatedResult[*entities.Category], error) {
	all := make([]*entities.Category, 0, len(m.categories))
	for _, c := range m.categories {
		all = append(all, c)
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
	return &service.PaginatedResult[*entities.Category]{
		Items:      all[start:end],
		Total:      total,
		Page:       params.Page,
		PerPage:    params.PerPage,
		TotalPages: totalPages,
	}, nil
}

func (m *mockCategoryRepo) UpdateCategory(c *entities.Category) error {
	m.categories[c.ID.String()] = c
	return nil
}

func (m *mockCategoryRepo) DeleteCategory(id uuid.UUID) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	if _, ok := m.categories[id.String()]; !ok {
		return service.ErrCategoryNotFound
	}
	delete(m.categories, id.String())
	return nil
}

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

func newTestCategoryHandler(repo *mockCategoryRepo) *CategoryHandler {
	return NewCategoryHandler(service.NewCategoryService(repo))
}

func categoryRouter(h *CategoryHandler) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/api/categories", h.CreateCategory).Methods("POST")
	r.HandleFunc("/api/categories", h.ListCategories).Methods("GET")
	r.HandleFunc("/api/categories/{id}", h.GetCategory).Methods("GET")
	r.HandleFunc("/api/categories/{id}", h.UpdateCategory).Methods("PUT")
	r.HandleFunc("/api/categories/{id}", h.DeleteCategory).Methods("DELETE")
	return r
}

func seedCategory(repo *mockCategoryRepo, name string) *entities.Category {
	c := entities.NewCategory(name)
	repo.categories[c.ID.String()] = c
	return c
}

// ---------------------------------------------------------------------------
// errorResponse
// ---------------------------------------------------------------------------

func TestErrorResponse(t *testing.T) {
	rr := httptest.NewRecorder()
	errorResponse(rr, "something went wrong", http.StatusBadRequest)

	assertStatus(t, rr, http.StatusBadRequest)
	assertContentType(t, rr, "application/json")
	assertJSONError(t, rr, "something went wrong")
}

// ---------------------------------------------------------------------------
// parsePaginationParams
// ---------------------------------------------------------------------------

func TestParsePaginationParams(t *testing.T) {
	t.Run("defaults when no query params", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/categories", nil)
		p := parsePaginationParams(req)
		if p.Page != 1 {
			t.Errorf("expected page 1, got %d", p.Page)
		}
		if p.PerPage != 10 {
			t.Errorf("expected per_page 10, got %d", p.PerPage)
		}
	})

	t.Run("custom values", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/categories?page=3&per_page=25", nil)
		p := parsePaginationParams(req)
		if p.Page != 3 {
			t.Errorf("expected page 3, got %d", p.Page)
		}
		if p.PerPage != 25 {
			t.Errorf("expected per_page 25, got %d", p.PerPage)
		}
	})

	t.Run("page below 1 defaults to 1", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/categories?page=0", nil)
		p := parsePaginationParams(req)
		if p.Page != 1 {
			t.Errorf("expected page 1, got %d", p.Page)
		}
	})

	t.Run("per_page above 100 defaults to 10", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/categories?per_page=200", nil)
		p := parsePaginationParams(req)
		if p.PerPage != 10 {
			t.Errorf("expected per_page 10, got %d", p.PerPage)
		}
	})

	t.Run("per_page below 1 defaults to 10", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/categories?per_page=-5", nil)
		p := parsePaginationParams(req)
		if p.PerPage != 10 {
			t.Errorf("expected per_page 10, got %d", p.PerPage)
		}
	})
}

// ---------------------------------------------------------------------------
// CreateCategory
// ---------------------------------------------------------------------------

func TestCreateCategory(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := newMockCategoryRepo()
		r := categoryRouter(newTestCategoryHandler(repo))

		body := bytes.NewBufferString(`{"name":"Electronics"}`)
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/categories", body)
		req.Header.Set("Content-Type", "application/json")

		r.ServeHTTP(rr, req)

		assertStatus(t, rr, http.StatusCreated)
		assertContentType(t, rr, "application/json")

		var resp entities.Category
		if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if resp.Name != "Electronics" {
			t.Errorf("expected name %q, got %q", "Electronics", resp.Name)
		}
		if resp.ID == uuid.Nil {
			t.Error("expected non-nil ID")
		}
		if len(repo.categories) != 1 {
			t.Errorf("expected 1 saved category, got %d", len(repo.categories))
		}
	})

	t.Run("invalid JSON body", func(t *testing.T) {
		repo := newMockCategoryRepo()
		r := categoryRouter(newTestCategoryHandler(repo))

		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/categories", bytes.NewBufferString(`{bad}`))
		r.ServeHTTP(rr, req)

		assertStatus(t, rr, http.StatusBadRequest)
		assertJSONError(t, rr, "invalid request body")
	})

	t.Run("empty name", func(t *testing.T) {
		repo := newMockCategoryRepo()
		r := categoryRouter(newTestCategoryHandler(repo))

		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/categories", bytes.NewBufferString(`{"name":""}`))
		r.ServeHTTP(rr, req)

		assertStatus(t, rr, http.StatusBadRequest)
		assertJSONError(t, rr, "name is required")
	})

	t.Run("whitespace-only name rejected by service", func(t *testing.T) {
		repo := newMockCategoryRepo()
		r := categoryRouter(newTestCategoryHandler(repo))

		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/categories", bytes.NewBufferString(`{"name":"   "}`))
		r.ServeHTTP(rr, req)

		assertStatus(t, rr, http.StatusBadRequest)
		assertJSONError(t, rr, "category name cannot be empty")
	})
}

// ---------------------------------------------------------------------------
// GetCategory
// ---------------------------------------------------------------------------

func TestGetCategory(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := newMockCategoryRepo()
		cat := seedCategory(repo, "Books")
		r := categoryRouter(newTestCategoryHandler(repo))

		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/categories/%s", cat.ID), nil)
		r.ServeHTTP(rr, req)

		assertStatus(t, rr, http.StatusOK)
		assertContentType(t, rr, "application/json")

		var resp entities.Category
		if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if resp.Name != "Books" {
			t.Errorf("expected name %q, got %q", "Books", resp.Name)
		}
	})

	t.Run("invalid ID", func(t *testing.T) {
		repo := newMockCategoryRepo()
		r := categoryRouter(newTestCategoryHandler(repo))

		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/categories/not-a-uuid", nil)
		r.ServeHTTP(rr, req)

		assertStatus(t, rr, http.StatusBadRequest)
		assertJSONError(t, rr, "invalid category id")
	})

	t.Run("not found", func(t *testing.T) {
		repo := newMockCategoryRepo()
		r := categoryRouter(newTestCategoryHandler(repo))

		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/categories/%s", uuid.New()), nil)
		r.ServeHTTP(rr, req)

		assertStatus(t, rr, http.StatusNotFound)
		assertJSONError(t, rr, "category not found")
	})
}

// ---------------------------------------------------------------------------
// ListCategories
// ---------------------------------------------------------------------------

func TestListCategories(t *testing.T) {
	t.Run("returns paginated result", func(t *testing.T) {
		repo := newMockCategoryRepo()
		seedCategory(repo, "Electronics")
		seedCategory(repo, "Books")
		r := categoryRouter(newTestCategoryHandler(repo))

		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/categories?page=1&per_page=10", nil)
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
		if resp.PerPage != 10 {
			t.Errorf("expected per_page 10, got %d", resp.PerPage)
		}
	})

	t.Run("empty list", func(t *testing.T) {
		repo := newMockCategoryRepo()
		r := categoryRouter(newTestCategoryHandler(repo))

		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/categories", nil)
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
// UpdateCategory
// ---------------------------------------------------------------------------

func TestUpdateCategory(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := newMockCategoryRepo()
		cat := seedCategory(repo, "Old Name")
		r := categoryRouter(newTestCategoryHandler(repo))

		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut,
			fmt.Sprintf("/api/categories/%s", cat.ID),
			bytes.NewBufferString(`{"name":"New Name"}`))
		r.ServeHTTP(rr, req)

		assertStatus(t, rr, http.StatusOK)
		assertContentType(t, rr, "application/json")

		var resp entities.Category
		if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if resp.Name != "New Name" {
			t.Errorf("expected name %q, got %q", "New Name", resp.Name)
		}
	})

	t.Run("invalid ID", func(t *testing.T) {
		repo := newMockCategoryRepo()
		r := categoryRouter(newTestCategoryHandler(repo))

		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, "/api/categories/bad-id",
			bytes.NewBufferString(`{"name":"Test"}`))
		r.ServeHTTP(rr, req)

		assertStatus(t, rr, http.StatusBadRequest)
		assertJSONError(t, rr, "invalid category id")
	})

	t.Run("invalid JSON body", func(t *testing.T) {
		repo := newMockCategoryRepo()
		cat := seedCategory(repo, "Test")
		r := categoryRouter(newTestCategoryHandler(repo))

		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut,
			fmt.Sprintf("/api/categories/%s", cat.ID),
			bytes.NewBufferString(`{bad}`))
		r.ServeHTTP(rr, req)

		assertStatus(t, rr, http.StatusBadRequest)
		assertJSONError(t, rr, "invalid request body")
	})

	t.Run("empty name", func(t *testing.T) {
		repo := newMockCategoryRepo()
		cat := seedCategory(repo, "Test")
		r := categoryRouter(newTestCategoryHandler(repo))

		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut,
			fmt.Sprintf("/api/categories/%s", cat.ID),
			bytes.NewBufferString(`{"name":""}`))
		r.ServeHTTP(rr, req)

		assertStatus(t, rr, http.StatusBadRequest)
		assertJSONError(t, rr, "name is required")
	})

	t.Run("not found", func(t *testing.T) {
		repo := newMockCategoryRepo()
		r := categoryRouter(newTestCategoryHandler(repo))

		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut,
			fmt.Sprintf("/api/categories/%s", uuid.New()),
			bytes.NewBufferString(`{"name":"New Name"}`))
		r.ServeHTTP(rr, req)

		assertStatus(t, rr, http.StatusNotFound)
		assertJSONError(t, rr, "category not found")
	})
}

// ---------------------------------------------------------------------------
// DeleteCategory
// ---------------------------------------------------------------------------

func TestDeleteCategory(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := newMockCategoryRepo()
		cat := seedCategory(repo, "To Delete")
		r := categoryRouter(newTestCategoryHandler(repo))

		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/categories/%s", cat.ID), nil)
		r.ServeHTTP(rr, req)

		assertStatus(t, rr, http.StatusNoContent)

		if _, ok := repo.categories[cat.ID.String()]; ok {
			t.Error("expected category to be deleted from repo")
		}
	})

	t.Run("invalid ID", func(t *testing.T) {
		repo := newMockCategoryRepo()
		r := categoryRouter(newTestCategoryHandler(repo))

		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/api/categories/bad-id", nil)
		r.ServeHTTP(rr, req)

		assertStatus(t, rr, http.StatusBadRequest)
		assertJSONError(t, rr, "invalid category id")
	})

	t.Run("not found", func(t *testing.T) {
		repo := newMockCategoryRepo()
		r := categoryRouter(newTestCategoryHandler(repo))

		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/categories/%s", uuid.New()), nil)
		r.ServeHTTP(rr, req)

		assertStatus(t, rr, http.StatusNotFound)
		assertJSONError(t, rr, "category not found")
	})
}
