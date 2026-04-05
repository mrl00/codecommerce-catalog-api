package service

import (
	"errors"
	"strings"

	"goapi/internal/entities"

	"github.com/google/uuid"
)

var ErrCategoryNotFound = errors.New("category not found")

type CategoryService struct {
	repo CategoryRepository
}

func (s *CategoryService) CreateCategory(name string) (*entities.Category, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("category name cannot be empty")
	}

	category := entities.NewCategory(name)
	if err := s.repo.SaveCategory(category); err != nil {
		return nil, err
	}

	return category, nil
}

func (s *CategoryService) GetCategory(id uuid.UUID) (*entities.Category, error) {
	category, err := s.repo.FindCategoryByID(id)
	if err != nil {
		return nil, ErrCategoryNotFound
	}
	return category, nil
}

func (s *CategoryService) ListCategories() ([]*entities.Category, error) {
	return s.repo.FindAllCategories()
}

func (s *CategoryService) UpdateCategory(id uuid.UUID, name string) (*entities.Category, error) {
	category, err := s.repo.FindCategoryByID(id)
	if err != nil {
		return nil, ErrCategoryNotFound
	}

	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("category name cannot be empty")
	}

	category.Name = name
	category.ResetUpdatedAt()

	if err := s.repo.UpdateCategory(category); err != nil {
		return nil, err
	}

	return category, nil
}

func (s *CategoryService) DeleteCategory(id uuid.UUID) error {
	_, err := s.repo.FindCategoryByID(id)
	if err != nil {
		return ErrCategoryNotFound
	}
	return s.repo.DeleteCategory(id)
}
