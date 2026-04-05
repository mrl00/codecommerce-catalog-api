package database

import (
	"database/sql"
	"codecommerceapi/internal/entities"
	"codecommerceapi/internal/service"

	"github.com/google/uuid"
)

type ProductDB struct {
	db *sql.DB
}

func NewProductDB(db *sql.DB) *ProductDB {
	return &ProductDB{db: db}
}

func (p *ProductDB) SaveProduct(product *entities.Product) error {
	_, err := p.db.Exec(`
		INSERT INTO tb_product (pk_product, tx_name, tx_description, nr_price, tx_image_url, fk_category, ts_product_created_at, ts_product_updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		product.ID, product.Name, product.Description, product.Price, product.ImageURL, product.CategoryID, product.CreatedAt, product.UpdatedAt)
	return err
}

func (p *ProductDB) FindProductByID(id uuid.UUID) (*entities.Product, error) {
	row := p.db.QueryRow(`
		SELECT pk_product, tx_name, tx_description, nr_price, tx_image_url, fk_category, ts_product_created_at, ts_product_updated_at
		FROM tb_product
		WHERE pk_product = $1`, id)

	product := &entities.Product{}
	err := row.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.ImageURL, &product.CategoryID, &product.CreatedAt, &product.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (p *ProductDB) FindAllProducts(params service.PaginationParams) (*service.PaginatedResult[*entities.Product], error) {
	var total int
	if err := p.db.QueryRow("SELECT COUNT(*) FROM tb_product").Scan(&total); err != nil {
		return nil, err
	}

	totalPages := (total + params.PerPage - 1) / params.PerPage
	offset := (params.Page - 1) * params.PerPage

	rows, err := p.db.Query(`
		SELECT pk_product, tx_name, tx_description, nr_price, tx_image_url, fk_category, ts_product_created_at, ts_product_updated_at
		FROM tb_product
		ORDER BY pk_product
		LIMIT $1 OFFSET $2`, params.PerPage, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*entities.Product
	for rows.Next() {
		product := &entities.Product{}
		err = rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.ImageURL, &product.CategoryID, &product.CreatedAt, &product.UpdatedAt)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &service.PaginatedResult[*entities.Product]{
		Items:      products,
		Total:      total,
		Page:       params.Page,
		PerPage:    params.PerPage,
		TotalPages: totalPages,
	}, nil
}

func (p *ProductDB) FindProductsByCategoryID(categoryID uuid.UUID, params service.PaginationParams) (*service.PaginatedResult[*entities.Product], error) {
	var total int
	if err := p.db.QueryRow("SELECT COUNT(*) FROM tb_product WHERE fk_category = $1", categoryID).Scan(&total); err != nil {
		return nil, err
	}

	totalPages := (total + params.PerPage - 1) / params.PerPage
	offset := (params.Page - 1) * params.PerPage

	rows, err := p.db.Query(`
		SELECT pk_product, tx_name, tx_description, nr_price, tx_image_url, fk_category, ts_product_created_at, ts_product_updated_at
		FROM tb_product
		WHERE fk_category = $1
		ORDER BY pk_product
		LIMIT $2 OFFSET $3`, categoryID, params.PerPage, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*entities.Product
	for rows.Next() {
		product := &entities.Product{}
		err = rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.ImageURL, &product.CategoryID, &product.CreatedAt, &product.UpdatedAt)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &service.PaginatedResult[*entities.Product]{
		Items:      products,
		Total:      total,
		Page:       params.Page,
		PerPage:    params.PerPage,
		TotalPages: totalPages,
	}, nil
}

func (p *ProductDB) UpdateProduct(product *entities.Product) error {
	_, err := p.db.Exec(`
		UPDATE tb_product
		SET tx_name = $2, tx_description = $3, nr_price = $4, tx_image_url = $5, fk_category = $6, ts_product_updated_at = $7
		WHERE pk_product = $1`,
		product.ID, product.Name, product.Description, product.Price, product.ImageURL, product.CategoryID, product.UpdatedAt)
	return err
}

func (p *ProductDB) DeleteProduct(id uuid.UUID) error {
	_, err := p.db.Exec("DELETE FROM tb_product WHERE pk_product = $1", id)
	return err
}
