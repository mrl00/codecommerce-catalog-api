package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"goapi/internal/database"
	"goapi/internal/entities"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// Category handlers

type CategoryHandler struct {
	db *database.CategoryDB
}

func NewCategoryHandler(db *database.CategoryDB) *CategoryHandler {
	return &CategoryHandler{db: db}
}

func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	category := entities.NewCategory(input.Name)
	if err := h.db.SaveCategory(category); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(category)
}

func (h *CategoryHandler) GetCategory(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "invalid category id", http.StatusBadRequest)
		return
	}

	category, err := h.db.FindCategoryByID(id)
	if err != nil {
		http.Error(w, "category not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

func (h *CategoryHandler) ListCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.db.FindAllCategories()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

func (h *CategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "invalid category id", http.StatusBadRequest)
		return
	}

	category, err := h.db.FindCategoryByID(id)
	if err != nil {
		http.Error(w, "category not found", http.StatusNotFound)
		return
	}

	var input struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	category.Name = input.Name
	category.UpdatedAt = time.Now()

	if err := h.db.UpdateCategory(category); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

func (h *CategoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "invalid category id", http.StatusBadRequest)
		return
	}

	if err := h.db.DeleteCategory(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Product handlers

type ProductHandler struct {
	db *database.ProductDB
}

func NewProductHandler(db *database.ProductDB) *ProductHandler {
	return &ProductHandler{db: db}
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name        string    `json:"name"`
		Description string    `json:"description"`
		Price       float64   `json:"price"`
		ImageURL    string    `json:"image_url"`
		CategoryID  uuid.UUID `json:"category"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	product := entities.NewProduct(input.Name, input.Description, input.Price, input.ImageURL, input.CategoryID)
	if err := h.db.SaveProduct(product); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}

func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "invalid product id", http.StatusBadRequest)
		return
	}

	product, err := h.db.FindProductByID(id)
	if err != nil {
		http.Error(w, "product not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

func (h *ProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.db.FindAllProducts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func (h *ProductHandler) ListProductsByCategory(w http.ResponseWriter, r *http.Request) {
	categoryID, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "invalid category id", http.StatusBadRequest)
		return
	}

	products, err := h.db.FindProductsByCategoryID(categoryID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "invalid product id", http.StatusBadRequest)
		return
	}

	product, err := h.db.FindProductByID(id)
	if err != nil {
		http.Error(w, "product not found", http.StatusNotFound)
		return
	}

	var input struct {
		Name        string    `json:"name"`
		Description string    `json:"description"`
		Price       float64   `json:"price"`
		ImageURL    string    `json:"image_url"`
		CategoryID  uuid.UUID `json:"category"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	product.Name = input.Name
	product.Description = input.Description
	product.Price = input.Price
	product.ImageURL = input.ImageURL
	product.CategoryID = input.CategoryID
	product.UpdatedAt = time.Now()

	if err := h.db.UpdateProduct(product); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "invalid product id", http.StatusBadRequest)
		return
	}

	if err := h.db.DeleteProduct(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
