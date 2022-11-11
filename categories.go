package main

import "time"

type Category struct {
	Id        int
	Name      string
	ImageUrl  string
	CreatedAt time.Time
}

type GetCategoryRes struct {
	Categories []*Category
	Count int
}

func (d *DBManager) CreateCategory(category *Category) (*Category, error) {
	query := `
		INSERT INTO categories(
			name,
			image_url
		) VALUES ($1, $2)
		RETURNING id, name, image_url, created_at
	`
	row := d.db.QueryRow(
		query,
		category.Name,
		category.ImageUrl,
	)
	var result Category
	err := row.Scan(
		&result.Id,
		&result.Name,
		&result.ImageUrl,
		&result.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (d *DBManager) GetCategory(category_id int) (*Category, error) {
	var result Category
	query := `
		SELECT
			id,
			name,
			image_url,
			created_at
		FROM categories
		WHERE id = $1
	`
	row := d.db.QueryRow(query, category_id)
	err := row.Scan(
		&result.Id,
		&result.Name,
		&result.ImageUrl,
		&result.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (d *DBManager) UpdateCategory(category *Category) (*Category, error) {
	query := `
		UPDATE categories 
		SET 
			name = $1,
			image_url = $2
		WHERE id = $3
		RETURNING id, name, image_url, created_at
	`
	row := d.db.QueryRow(
		query,
		category.Name,
		category.ImageUrl,
		category.Id,
	)
	var result Category
	err := row.Scan(
		&result.Id,
		&result.Name,
		&result.ImageUrl,
		&result.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (d *DBManager) DeleteCategory(id int) error {
	query := `
		DELETE FROM categories WHERE id = $1  
	`
	_, err := d.db.Exec(query, id)
	return err
}

func (d *DBManager) GetAllCategories() (*GetCategoryRes, error) {
	var result GetCategoryRes
	result.Categories = make([]*Category, 0)
	query := `
		SELECT 
			id, 
			name, 
			image_url,
			created_at
		FROM categories
	`
	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var category Category
		err := rows.Scan(
			&category.Id,
			&category.Name,
			&category.ImageUrl,
			&category.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		result.Categories = append(result.Categories, &category)
	}
	queryCount := "SELECT count(*) FROM categories"
	err = d.db.QueryRow(queryCount).Scan(&result.Count)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
