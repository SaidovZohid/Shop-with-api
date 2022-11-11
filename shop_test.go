package main

import (
	"testing"

	_ "github.com/bxcodec/faker/v4"
	"github.com/stretchr/testify/require"
)

// func createCustomer(t *testing.T) *Customer {
// 	return nil
// }

// func createCategory(t *testing.T) *Customer {
// 	return nil
// }

// func createProduct(t *testing.T) *Customer {
// 	return nil
// }

func createOrder(t *testing.T) int {
	order_id, err := dbManager.CreateOrder(&Orders{
		CustomerId: 1001,
		TotalAmount: 9000,
		Items: []*OrderItem{
			{
				ProductName: "Iphone 13 Pro Max",
				ProductId: 4,
				Count: 2,
				TotalPrice: 2000,
				Status: true,
			},
			{
				ProductName: "Iphone 11 Pro Max",
				ProductId: 2,
				Count: 4,
				TotalPrice: 4000,
				Status: true,
			},
			{
				ProductName: "Iphone Se",
				ProductId: 3,
				Count: 3,
				TotalPrice: 3000,
				Status: true,
			},
		},
	}) 
	require.NoError(t, err)
	require.NotEmpty(t, order_id)
	return order_id
}

func deleteOrder(t *testing.T, order_id int) {
	dbManager.DeleteOrder(order_id)
}

func TestCreateOrder(t *testing.T) {
	order_id := createOrder(t)
	require.NotEmpty(t, order_id)
	deleteOrder(t, order_id)
}

func TestGetOrder(t *testing.T) {
	order_id := createOrder(t)
	orderitem, err := dbManager.GetOrder(order_id)
	require.NoError(t, err)
	require.NotEmpty(t, orderitem)
	deleteOrder(t, order_id)
}

func TestUpdateOrder(t *testing.T) {
	order_id := createOrder(t)
	err := dbManager.UpdateOrder(&Orders{
		CustomerId: 1001,
		TotalAmount: 13000,
		Items: []*OrderItem{
			{
				ProductName: "Iphone 11 Pro Max",
				ProductId: 2,
				Count: 2,
				TotalPrice: 6000,
				Status: true,
			},
			{
				ProductName: "Iphone Se",
				ProductId: 3,
				Count: 4,
				TotalPrice: 7000,
				Status: true,
			},
		},
	})
	require.NoError(t, err)
	require.NotEmpty(t, order_id)
	deleteOrder(t, order_id)
}

func TestDeleteOrder(t *testing.T) {
	order_id := createOrder(t)
	err := dbManager.DeleteOrder(order_id)
	require.NoError(t, err)
	require.NotEmpty(t, order_id)
}

func TestGetAllOrder(t *testing.T) {
	order_id := createOrder(t)
	orders, err := dbManager.GetAllOrders(&GetAllOrderss{
		Limit: 10,
		Page: 1,
		CustomerId: 1001,
	})
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(orders.Orders), 1)
	deleteOrder(t, order_id)
}