package main

import (
	"database/sql"
	"fmt"
	"time"
)

type GetProducts struct {
	Products []*Product
}

type DBManager struct {
	db *sql.DB
}

type Product struct {
	Id            int
	Category_id   int
	Category_name string
	Name          string
	Price         float64
	ImageUrl      string
	CreatedAt     time.Time
	Images        []*ProductImages
}

type ProductImages struct {
	Id             int
	ImageUrl       string
	SequenceNumber int
}

type GetProductsParams struct {
	Limit       int    `json:"Limit"`
	Page        int    `json:"Page"`
	ProductName string `json:"product_name"`
}

func (d *DBManager) CreateProduct(product *Product) (int64, error) {
	var productId int64
	query := `
		INSERT INTO products (
			category_id,
			name,
			price,
			image_url
		) VALUES ($1, $2, $3, $4) RETURNING id
	`
	row := d.db.QueryRow(
		query,
		product.Category_id,
		product.Name,
		product.Price,
		product.ImageUrl,
	)
	err := row.Scan(&productId)
	if err != nil {
		return 0, err
	}
	queryInsertImages := `
		INSERT INTO product_images(
			product_id,
			image_url,
			sequence_number
		) VALUES ($1, $2, $3)
	`
	for _, image := range product.Images {
		_, err := d.db.Exec(
			queryInsertImages,
			productId,
			image.ImageUrl,
			image.SequenceNumber,
		)
		if err != nil {
			return 0, err
		}
	}
	return productId, nil
}

func (d *DBManager) GetProduct(id int64) (*Product, error) {
	var product Product
	product.Images = make([]*ProductImages, 0)
	query := `
		SELECT 
			p.id,
			p.category_id,
			c.name,
			p.name,
			p.price,
			p.image_url,
			p.created_at
		FROM products p INNER JOIN categories c on c.id = p.category_id 
		WHERE p.id = $1
	`
	row := d.db.QueryRow(query, id)
	err := row.Scan(
		&product.Id,
		&product.Category_id,
		&product.Category_name,
		&product.Name,
		&product.Price,
		&product.ImageUrl,
		&product.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	queryImages := `
		SELECT 
			id, 
			image_url,
			sequence_number
		FROM product_images 
		WHERE product_id = $1
	`
	rows, err := d.db.Query(queryImages, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var image ProductImages
		err := rows.Scan(
			&image.Id,
			&image.ImageUrl,
			&image.SequenceNumber,
		)
		if err != nil {
			return nil, err
		}
		product.Images = append(product.Images, &image)
	}
	return &product, nil
}

func (d *DBManager) GetAllProducts(params *GetProductsParams) (*GetProducts, error) {
	var result GetProducts
	result.Products = make([]*Product, 0)
	filter := ""
	if params.ProductName != "" {
		filter = fmt.Sprintf("WHERE p.name ilike '%s'", "%"+params.ProductName+"%")
	}
	query := `
		SELECT 
			p.id,
			p.category_id,
			c.name,
			p.name,
			p.price,
			p.image_url,
			p.created_at
		FROM products p 
		INNER JOIN categories c on c.id = p.category_id
	` + filter + `
		ORDER BY created_at DESC 
		LIMIT $1 OFFSET $2
	`
	offset := (params.Page - 1) * params.Limit
	rows, err := d.db.Query(query, params.Limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var product Product
		err := rows.Scan(
			&product.Id,
			&product.Category_id,
			&product.Category_name,
			&product.Name,
			&product.Price,
			&product.ImageUrl,
			&product.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		result.Products = append(result.Products, &product)
	}
	return &result, nil
}

func (d *DBManager) UpdateProduct(product *Product) error {
	query := `
		UPDATE products SET 
			category_id = $1,
			name = $2,
			price = $3,
			image_url = $4
		WHERE id = $5
	`
	result, err := d.db.Exec(
		query,
		product.Category_id,
		product.Name,
		product.Price,
		product.ImageUrl,
		product.Id,
	)
	if err != nil {
		return err
	}
	rowsCount, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsCount == 0 {
		return sql.ErrNoRows
	}
	queryDeleteImages := `DELETE FROM product_images WHERE product_id = $1`
	_, err = d.db.Exec(queryDeleteImages, product.Id)
	if err != nil {
		return err
	}
	queryInsertImages := `
		INSERT INTO product_images (
			image_url,
			sequence_number,
			product_id
		) VALUES ($1, $2, $3)
	`
	for _, image := range product.Images {
		_, err = d.db.Exec(
			queryInsertImages,
			image.ImageUrl,
			image.SequenceNumber,
			product.Id,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *DBManager) DeleteProduct(id int64) error {
	queryDeleteImages := `DELETE FROM product_images WHERE product_id = $1`
	_, err := d.db.Exec(queryDeleteImages, id)
	if err != nil {
		return err
	}
	queryDeleteProduct := `DELETE FROM products WHERE id = $1`
	result, err := d.db.Exec(queryDeleteProduct, id)
	if err != nil {
		return err
	}
	rowsCount, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsCount == 0 {
		return sql.ErrNoRows
	}
	return nil
}
