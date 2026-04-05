package service

import (
	"errors"
	"strings"

	"codecommerceapi/internal/entities"

	"github.com/google/uuid"
)

type ProductService struct {
	repo ProductRepository
}

func (s *ProductService) CreateProduct(name, description string, price int64, imageURL string, categoryID uuid.UUID) (*entities.Product, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("product name cannot be empty")
	}

	if price < 0 {
		return nil, errors.New("product price cannot be negative")
	}

	if categoryID == uuid.Nil {
		return nil, errors.New("product category ID is required")
	}

	product := entities.NewProduct(name, description, price, imageURL, categoryID)
	if err := s.repo.SaveProduct(product); err != nil {
		return nil, err
	}

	return product, nil
}

func (s *ProductService) GetProduct(id uuid.UUID) (*entities.Product, error) {
	product, err := s.repo.FindProductByID(id)
	if err != nil {
		return nil, ErrProductNotFound
	}
	return product, nil
}

func (s *ProductService) ListProducts(params PaginationParams) (*PaginatedResult[*entities.Product], error) {
	return s.repo.FindAllProducts(params)
}

func (s *ProductService) ListProductsByCategory(categoryID uuid.UUID, params PaginationParams) (*PaginatedResult[*entities.Product], error) {
	return s.repo.FindProductsByCategoryID(categoryID, params)
}

func (s *ProductService) UpdateProduct(id uuid.UUID, name, description string, price int64, imageURL string, categoryID uuid.UUID) (*entities.Product, error) {
	product, err := s.repo.FindProductByID(id)
	if err != nil {
		return nil, ErrProductNotFound
	}

	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("product name cannot be empty")
	}

	if price < 0 {
		return nil, errors.New("product price cannot be negative")
	}

	if categoryID == uuid.Nil {
		return nil, errors.New("product category ID is required")
	}

	product.Name = name
	product.Description = description
	product.Price = price
	product.ImageURL = imageURL
	product.CategoryID = categoryID
	product.ResetUpdatedAt()

	if err := s.repo.UpdateProduct(product); err != nil {
		return nil, err
	}

	return product, nil
}

func (s *ProductService) DeleteProduct(id uuid.UUID) error {
	return s.repo.DeleteProduct(id)
}
