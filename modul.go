package main

import (
	"time"
)

type ResponseError struct {
	Message string `json:"message"`
}

type ResponseOk struct {
	Message string `json:"message"`
}

type CreateCategory struct {
	Name      string `json:"name"`
	ImageUrl  string `json:""`
}

type CreateorUpdateorGetCustomer struct {
	Id          int       `json:"id"`
	FirstName   string    `json:"firstname"`
	LastName    string    `json:"lastname"`
	PhoneNumber string    `json:"phone_number"`
	Gender      bool      `json:"gender"`
	BirthDate   time.Time `json:"birth_date"`
	Balance     float64   `json:"balance"`
}