package database

import (
	"database/sql"
	"codecommerceapi/internal/entities"

	"github.com/google/uuid"
)

type CategoryDB struct {
	db *sql.DB
}

func NewCategoryDB(db *sql.DB) *CategoryDB {
	return &CategoryDB{db: db}
}

func (c *CategoryDB) SaveCategory(category *entities.Category) error {
	stmt, err := c.db.Prepare(`
		INSERT INTO tb_category (pk_category, tx_name, ts_category_created_at, ts_category_updated_at)
		VALUES ($1, $2, $3, $4)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(category.ID, category.Name, category.CreatedAt, category.UpdatedAt)
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

func (c *CategoryDB) FindAllCategories() ([]*entities.Category, error) {
	rows, err := c.db.Query(`
		SELECT pk_category, tx_name, ts_category_created_at, ts_category_updated_at
		FROM tb_category`)
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
	return categories, nil
}

func (c *CategoryDB) UpdateCategory(category *entities.Category) error {
	stmt, err := c.db.Prepare(`
		UPDATE tb_category
		SET tx_name = $2, ts_category_updated_at = $3
		WHERE pk_category = $1`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(category.ID, category.Name, category.UpdatedAt)
	return err
}

func (c *CategoryDB) DeleteCategory(id uuid.UUID) error {
	stmt, err := c.db.Prepare("DELETE FROM tb_category WHERE pk_category = $1")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	return err
}
