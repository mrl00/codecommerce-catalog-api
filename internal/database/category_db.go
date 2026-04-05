package database

import (
	"database/sql"
	"codecommerceapi/internal/entities"
	"codecommerceapi/internal/service"

	"github.com/google/uuid"
)

type CategoryDB struct {
	db *sql.DB
}

func NewCategoryDB(db *sql.DB) *CategoryDB {
	return &CategoryDB{db: db}
}

func (c *CategoryDB) SaveCategory(category *entities.Category) error {
	_, err := c.db.Exec(`
		INSERT INTO tb_category (pk_category, tx_name, ts_category_created_at, ts_category_updated_at)
		VALUES ($1, $2, $3, $4)`,
		category.ID, category.Name, category.CreatedAt, category.UpdatedAt)
	return err
}

func (c *CategoryDB) FindCategoryByID(id uuid.UUID) (*entities.Category, error) {
	row := c.db.QueryRow(`
		SELECT pk_category, tx_name, ts_category_created_at, ts_category_updated_at
		FROM tb_category
		WHERE pk_category = $1`, id)

	category := &entities.Category{}
	err := row.Scan(&category.ID, &category.Name, &category.CreatedAt, &category.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (c *CategoryDB) FindAllCategories(params service.PaginationParams) (*service.PaginatedResult[*entities.Category], error) {
	var total int
	if err := c.db.QueryRow("SELECT COUNT(*) FROM tb_category").Scan(&total); err != nil {
		return nil, err
	}

	totalPages := (total + params.PerPage - 1) / params.PerPage
	offset := (params.Page - 1) * params.PerPage

	rows, err := c.db.Query(`
		SELECT pk_category, tx_name, ts_category_created_at, ts_category_updated_at
		FROM tb_category
		ORDER BY pk_category
		LIMIT $1 OFFSET $2`, params.PerPage, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*entities.Category
	for rows.Next() {
		category := &entities.Category{}
		err = rows.Scan(&category.ID, &category.Name, &category.CreatedAt, &category.UpdatedAt)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &service.PaginatedResult[*entities.Category]{
		Items:      categories,
		Total:      total,
		Page:       params.Page,
		PerPage:    params.PerPage,
		TotalPages: totalPages,
	}, nil
}

func (c *CategoryDB) UpdateCategory(category *entities.Category) error {
	_, err := c.db.Exec(`
		UPDATE tb_category
		SET tx_name = $2, ts_category_updated_at = $3
		WHERE pk_category = $1`,
		category.ID, category.Name, category.UpdatedAt)
	return err
}

func (c *CategoryDB) DeleteCategory(id uuid.UUID) error {
	result, err := c.db.Exec("DELETE FROM tb_category WHERE pk_category = $1", id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return service.ErrCategoryNotFound
	}
	return nil
}
