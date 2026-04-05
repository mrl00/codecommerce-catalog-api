package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"codecommerceapi/internal/service"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type ProductHandler struct {
	svc *service.ProductService
}

func NewProductHandler(svc *service.ProductService) *ProductHandler {
	return &ProductHandler{svc: svc}
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Price       int64  `json:"price"`
		ImageURL    string `json:"image_url"`
		CategoryID  string `json:"category"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		errorResponse(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if input.Name == "" {
		errorResponse(w, "name is required", http.StatusBadRequest)
		return
	}

	if input.Price <= 0 {
		errorResponse(w, "price must be greater than zero", http.StatusBadRequest)
		return
	}

	if input.CategoryID == "" {
		errorResponse(w, "category is required", http.StatusBadRequest)
		return
	}

	categoryID, err := uuid.Parse(input.CategoryID)
	if err != nil {
		errorResponse(w, "invalid category id format", http.StatusBadRequest)
		return
	}

	product, err := h.svc.CreateProduct(input.Name, input.Description, input.Price, input.ImageURL, categoryID)
	if err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}

func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		errorResponse(w, "invalid product id", http.StatusBadRequest)
		return
	}

	product, err := h.svc.GetProduct(id)
	if err != nil {
		if err == service.ErrProductNotFound {
			errorResponse(w, "product not found", http.StatusNotFound)
			return
		}
		errorResponse(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

func (h *ProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.svc.ListProducts()
	if err != nil {
		errorResponse(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func (h *ProductHandler) ListProductsByCategory(w http.ResponseWriter, r *http.Request) {
	categoryID, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		errorResponse(w, "invalid category id", http.StatusBadRequest)
		return
	}

	products, err := h.svc.ListProductsByCategory(categoryID)
	if err != nil {
		errorResponse(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		errorResponse(w, "invalid product id", http.StatusBadRequest)
		return
	}

	var input struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Price       int64  `json:"price"`
		ImageURL    string `json:"image_url"`
		CategoryID  string `json:"category"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		errorResponse(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if input.Name == "" {
		errorResponse(w, "name is required", http.StatusBadRequest)
		return
	}

	if input.Price <= 0 {
		errorResponse(w, "price must be greater than zero", http.StatusBadRequest)
		return
	}

	if input.CategoryID == "" {
		errorResponse(w, "category is required", http.StatusBadRequest)
		return
	}

	categoryID, err := uuid.Parse(input.CategoryID)
	if err != nil {
		errorResponse(w, "invalid category id format", http.StatusBadRequest)
		return
	}

	product, err := h.svc.UpdateProduct(id, input.Name, input.Description, input.Price, input.ImageURL, categoryID)
	if err != nil {
		if err == service.ErrProductNotFound {
			errorResponse(w, "product not found", http.StatusNotFound)
			return
		}
		errorResponse(w, fmt.Sprintf("failed to update product: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		errorResponse(w, "invalid product id", http.StatusBadRequest)
		return
	}

	if err := h.svc.DeleteProduct(id); err != nil {
		if err == service.ErrProductNotFound {
			errorResponse(w, "product not found", http.StatusNotFound)
			return
		}
		errorResponse(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
