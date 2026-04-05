package service

import (
	"codecommerceapi/internal/entities"

	"github.com/google/uuid"
)

type CategoryRepository interface {
	SaveCategory(category *entities.Category) error
	FindCategoryByID(id uuid.UUID) (*entities.Category, error)
	FindAllCategories() ([]*entities.Category, error)
	UpdateCategory(category *entities.Category) error
	DeleteCategory(id uuid.UUID) error
}

type ProductRepository interface {
	SaveProduct(product *entities.Product) error
	FindProductByID(id uuid.UUID) (*entities.Product, error)
	FindAllProducts() ([]*entities.Product, error)
	FindProductsByCategoryID(categoryID uuid.UUID) ([]*entities.Product, error)
	UpdateProduct(product *entities.Product) error
	DeleteProduct(id uuid.UUID) error
}

func NewCategoryService(repo CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func NewProductService(repo ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}
