package main

import (
	"fmt"
	"time"
)

type Customer struct {
	Id          int       `json:"id"`
	FirstName   string    `json:"firstname"`
	LastName    string    `json:"lastname"`
	PhoneNumber string    `json:"phone_number"`
	Gender      bool      `json:"gender"`
	BirthDate   time.Time `json:"birth_date"`
	Balance     float64   `json:"balance"`
	Created_at  time.Time `json:"created_at"`
	Updated_at  time.Time `json:"updated_at"`
	Deleted_at  time.Time `json:"deleted_at"`
}

type GetCustomers struct {
	Customers []*Customer `json:"customers"`
	Count     int         `json:"count"`
}

type GetCustomersParams struct {
	Limit        int    `json:"limit"`
	Page         int    `json:"page"`
	CustomerName string `json:"name"`
}

func (d *DBManager) CreateCustomer(customer *CreateorUpdateorGetCustomer) (*Customer, error) {
	var result Customer
	query := `
		INSERT INTO customer(
			firstname,
			lastname,
			phone_number,
			gender,
			birth_date
		) VALUES ($1, $2, $3, $4, $5)
		RETURNING id, firstname, lastname, phone_number, gender, birth_date, created_at
	`
	row := d.db.QueryRow(
		query,
		customer.FirstName,
		customer.LastName,
		customer.PhoneNumber,
		customer.Gender,
		customer.BirthDate,
	)
	err := row.Scan(
		&result.Id,
		&result.FirstName,
		&result.LastName,
		&result.PhoneNumber,
		&result.Gender,
		&result.BirthDate,
		&result.Created_at,
	)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (d *DBManager) GetCustomer(customer_id int) (*CreateorUpdateorGetCustomer, error) {
	var result CreateorUpdateorGetCustomer
	query := `
		SELECT 
			id,
			firstname,
			lastname,
			phone_number,
			gender,
			birth_date
		FROM customer WHERE id = $1
	`
	row := d.db.QueryRow(query, customer_id)
	err := row.Scan(
		&result.Id,
		&result.FirstName,
		&result.LastName,
		&result.PhoneNumber,
		&result.Gender,
		&result.BirthDate,
	)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (d *DBManager) UpdateCustomer(customer *CreateorUpdateorGetCustomer) (*CreateorUpdateorGetCustomer, error) {
	var result CreateorUpdateorGetCustomer
	query := `
		UPDATE customer 
		SET 
			firstname = $1,
			lastname = $2,
			phone_number = $3,
			gender = $4,
			birth_date = $5,
			updated_at = $6
		WHERE id = $7
		RETURNING id, firstname, lastname, phone_number, gender, birth_date
	`
	row := d.db.QueryRow(
		query,
		customer.FirstName,
		customer.LastName,
		customer.PhoneNumber,
		customer.Gender,
		customer.BirthDate,
		time.Now(),
		customer.Id,
	)
	err := row.Scan(
		&result.Id,
		&result.FirstName,
		&result.LastName,
		&result.PhoneNumber,
		&result.Gender,
		&result.BirthDate,
	)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (d *DBManager) DeleteCustomer(customer_id int) error {
	query := `
		UPDATE customer SET 
			deleted_at = $1
		WHERE id = $2
	`
	_, err := d.db.Exec(query, time.Now(), customer_id)
	return err
}

func (d *DBManager) GetAllCustomers(params *GetCustomersParams) (*GetCustomers, error) {
	var result GetCustomers
	result.Customers = make([]*Customer, 0)
	offset := (params.Page - 1) * params.Limit
	filter := ""
	if params.CustomerName != "" {
		filter = fmt.Sprintf("WHERE name ilike '%s'", "%"+params.CustomerName+"%")
	}
	query := `
		SELECT 
			id, 
			firstname,
			lastname,
			phone_number,
			gender,
			birth_date
		FROM customer 
	` + filter + `LIMIT $1 OFFSET $2`
	rows, err := d.db.Query(query, params.Limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var c Customer
		err := rows.Scan(
			&c.Id,
			&c.FirstName,
			&c.LastName,
			&c.PhoneNumber,
			&c.Gender,
			&c.BirthDate,
		)
		if err != nil {
			return nil, err
		}
		result.Customers = append(result.Customers, &c)
	}
	queryCount := "SELECT count(*) from customer " + filter
	err = d.db.QueryRow(queryCount).Scan(&result.Count)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
