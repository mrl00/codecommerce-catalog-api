package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"codecommerceapi/internal/service"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type CategoryHandler struct {
	svc *service.CategoryService
}

func NewCategoryHandler(svc *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{svc: svc}
}

func parseCategoryPaginationParams(r *http.Request) service.PaginationParams {
	q := r.URL.Query()
	page, perPage := 1, 10
	if v := q.Get("page"); v != "" {
		fmt.Sscanf(v, "%d", &page)
	}
	if v := q.Get("per_page"); v != "" {
		fmt.Sscanf(v, "%d", &perPage)
	}
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}
	return service.PaginationParams{Page: page, PerPage: perPage}
}

func errorResponse(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		errorResponse(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if input.Name == "" {
		errorResponse(w, "name is required", http.StatusBadRequest)
		return
	}

	category, err := h.svc.CreateCategory(input.Name)
	if err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(category)
}

func (h *CategoryHandler) GetCategory(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		errorResponse(w, "invalid category id", http.StatusBadRequest)
		return
	}

	category, err := h.svc.GetCategory(id)
	if err != nil {
		if err == service.ErrCategoryNotFound {
			errorResponse(w, "category not found", http.StatusNotFound)
			return
		}
		errorResponse(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

func (h *CategoryHandler) ListCategories(w http.ResponseWriter, r *http.Request) {
	params := parseCategoryPaginationParams(r)
	result, err := h.svc.ListCategories(params)
	if err != nil {
		errorResponse(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (h *CategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		errorResponse(w, "invalid category id", http.StatusBadRequest)
		return
	}

	var input struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		errorResponse(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if input.Name == "" {
		errorResponse(w, "name is required", http.StatusBadRequest)
		return
	}

	category, err := h.svc.UpdateCategory(id, input.Name)
	if err != nil {
		if err == service.ErrCategoryNotFound {
			errorResponse(w, "category not found", http.StatusNotFound)
			return
		}
		errorResponse(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

func (h *CategoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		errorResponse(w, "invalid category id", http.StatusBadRequest)
		return
	}

	if err := h.svc.DeleteCategory(id); err != nil {
		if err == service.ErrCategoryNotFound {
			errorResponse(w, "category not found", http.StatusNotFound)
			return
		}
		errorResponse(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
