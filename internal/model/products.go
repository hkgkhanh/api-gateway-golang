package model

// User is the db schema for the users table
type Product struct {
	ID           int    `json:"id" db:"id"`
	Name         string `json:"name" db:"name"`
	Manufacturer string `json:"manufacturer" db:"manufacturer"`
	Quantity     int    `json:"quantity" db:"quantity"`
	Base
}

// Below are the structures of the request/response structs in the users handler.

type AddProductRequest struct {
	Name         string `json:"name" validate:"required"`
	Manufacturer string `json:"manufacturer" validate:"required"`
	Quantity     int    `json:"quantity" validate:"required"`
}

type GetProductResponse struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Manufacturer string `json:"manufacturer"`
	Quantity     int    `json:"quantity"`
}

type UpdateProductRequest struct {
	ID           int    `json:"id" validate:"required"`
	Name         string `json:"name"`
	Manufacturer string `json:"manufacturer"`
	Quantity     int    `json:"quantity"`
}

// type IDResponse struct {
// 	ID int `json:"id"`
// }
