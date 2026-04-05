package handler

import (
	"encoding/json"
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
	categories, err := h.svc.ListCategories()
	if err != nil {
		errorResponse(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
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
