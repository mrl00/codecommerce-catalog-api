package service

import (
	"errors"

	"codecommerceapi/internal/entities"

	"github.com/google/uuid"
)

var (
	ErrCategoryNotFound = errors.New("category not found")
	ErrProductNotFound  = errors.New("product not found")
)

type PaginationParams struct {
	Page, PerPage int
}

type PaginatedResult[T any] struct {
	Items      []T `json:"items"`
	Total      int `json:"total"`
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	TotalPages int `json:"total_pages"`
}

type CategoryRepository interface {
	SaveCategory(category *entities.Category) error
	FindCategoryByID(id uuid.UUID) (*entities.Category, error)
	FindAllCategories(params PaginationParams) (*PaginatedResult[*entities.Category], error)
	UpdateCategory(category *entities.Category) error
	DeleteCategory(id uuid.UUID) error
}

type ProductRepository interface {
	SaveProduct(product *entities.Product) error
	FindProductByID(id uuid.UUID) (*entities.Product, error)
	FindAllProducts(params PaginationParams) (*PaginatedResult[*entities.Product], error)
	FindProductsByCategoryID(categoryID uuid.UUID, params PaginationParams) (*PaginatedResult[*entities.Product], error)
	UpdateProduct(product *entities.Product) error
	DeleteProduct(id uuid.UUID) error
}

func NewCategoryService(repo CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func NewProductService(repo ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}
