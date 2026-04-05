package database

import (
	"database/sql"
	"codecommerceapi/internal/entities"

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

func (p *ProductDB) FindAllProducts() ([]*entities.Product, error) {
	rows, err := p.db.Query(`
		SELECT pk_product, tx_name, tx_description, nr_price, tx_image_url, fk_category, ts_product_created_at, ts_product_updated_at
		FROM tb_product`)
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
	return products, nil
}

func (p *ProductDB) FindProductsByCategoryID(categoryID uuid.UUID) ([]*entities.Product, error) {
	rows, err := p.db.Query(`
		SELECT pk_product, tx_name, tx_description, nr_price, tx_image_url, fk_category, ts_product_created_at, ts_product_updated_at
		FROM tb_product
		WHERE fk_category = $1`, categoryID)
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
	return products, nil
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
	stmt, err := p.db.Prepare("DELETE FROM tb_product WHERE pk_product = $1")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	return err
}
