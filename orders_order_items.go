package main

import (
	"database/sql"
	"fmt"
)

type Orders struct {
	Id          int
	CustomerId  int
	Items       []*OrderItem
	TotalAmount float64
}

type OrderItem struct {
	Id          int
	OrderId     int
	ProductName string
	ProductId   int
	Count       int
	TotalPrice  float64
	Status      bool
}

type GetAllOrders struct {
	Orders []*Orders `json:"orders"`
	Count  int       `json:"count"`
}

type GetAllParams struct {
	Limit      int `json:"limit"`
	Page       int `json:"page"`
	CustomerID int `json:"customer_id"`
}

type GetAllOrderss struct {
	Limit      int 
	Page       int
	CustomerId int
}

func (d *DBManager) CreateOrder(o *Orders) (int, error) {
	tx, err := d.db.Begin()
	if err != nil {
		return 0, err
	}
	queryUpdateBalance := `UPDATE customer SET balance = balance - $1 WHERE id = $2`
	_, err = tx.Exec(queryUpdateBalance, o.TotalAmount, o.CustomerId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	var order_id int
	query := `
		INSERT INTO orders(
			customer_id,
			total_amount
		) VALUES ($1, $2)
		RETURNING customer_id, id
	`
	row := tx.QueryRow(
		query,
		o.CustomerId,
		o.TotalAmount,
	)
	var id int
	err = row.Scan(
		&id,
		&order_id,
	)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	queryOrderItems := `
		INSERT INTO order_items(
			order_id,
			product_name,
			product_id,
			count,
			total_price,
			status
		) VALUES ($1, $2, $3, $4, $5, $6)
	`
	for _, v := range o.Items {
		_, err := tx.Exec(
			queryOrderItems,
			order_id,
			v.ProductName,
			v.ProductId,
			v.Count,
			v.TotalPrice,
			v.Status,
		)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}
	err = tx.Commit()
	if err != nil {
		return 0, err
	}
	return order_id, nil
}

func (d *DBManager) GetOrder(order_id int) ([]*OrderItem, error) {
	var result []*OrderItem
	queryOrderItems := `
		SELECT 
			o.id,
			o.order_id,
			o.product_name,
			o.product_id,
			o.count,
			o.total_price,
			o.status
		FROM order_items o WHERE o.order_id = $1
	`
	rows, err := d.db.Query(queryOrderItems, order_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var order_items OrderItem
		err := rows.Scan(
			&order_items.Id,
			&order_items.OrderId,
			&order_items.ProductName,
			&order_items.ProductId,
			&order_items.Count,
			&order_items.TotalPrice,
			&order_items.Status,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, &order_items)
	}
	return result, nil
}

func (d *DBManager) UpdateOrder(orders *Orders) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	query := `
		UPDATE orders SET
			total_amount = $1
		WHERE customer_id = $2
		RETURNING id
	`
	result := tx.QueryRow(
		query,
		orders.TotalAmount,
		orders.CustomerId,
	)
	var order_id int
	err = result.Scan(
		&order_id,
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	query = `DELETE FROM order_items WHERE order_id = $1`
	_, err = tx.Exec(query, order_id)
	if err != nil {
		tx.Rollback()
		return err
	}
	query = `
		INSERT INTO order_items(
			order_id,
			product_name,
			product_id,
			count,
			total_price,
			status
		) VALUES ($1, $2, $3, $4, $5, $6)
	`
	for _, v := range orders.Items {
		_, err = tx.Exec(
			query,
			order_id,
			v.ProductName,
			v.ProductId,
			v.Count,
			v.TotalPrice,
			v.Status,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (d *DBManager) DeleteOrder(order_id int) error {
	query := `DELETE FROM order_items WHERE order_id = $1`
	row, err := d.db.Exec(query, order_id)
	if err != nil {
		return err
	}
	res, err := row.RowsAffected()
	if err != nil {
		return err
	}
	if res == 0 {
		return sql.ErrNoRows
	}
	queryDeleteOrders := `DELETE FROM orders WHERE id = $1`
	_, err = d.db.Exec(queryDeleteOrders, order_id)
	if err != nil {
		return err
	}
	return nil
}

func (d *DBManager) GetAllOrders(params *GetAllOrderss) (*GetAllOrders, error) {
	var result GetAllOrders
	var order Orders
	result.Orders = make([]*Orders, 0)
	filter := ""
	offset := (params.Page - 1) * params.Limit
	if params.CustomerId != 0 {
		filter = fmt.Sprintf(" WHERE o.customer_id = %d", params.CustomerId)
	}
	query := `
		SELECT 
			id,
			customer_id,
			total_amount
		FROM orders o
	` + filter + ` LIMIT $1 OFFSET $2`
	rows, err := d.db.Query(query, params.Limit, offset)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err := rows.Scan(
			&order.Id,
			&order.CustomerId,
			&order.TotalAmount,
		)
		if err != nil {
			return nil, err
		}
		result.Orders = append(result.Orders, &order)
	}
	return &result, nil
}
