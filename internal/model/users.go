package model

// User is the db schema for the users table
type User struct {
	ID          int    `json:"id" db:"id"`
	FirstName   string `json:"firstName" db:"first_name"`
	LastName    string `json:"lastName" db:"last_name"`
	Email       string `json:"email" db:"email"`
	PhoneNumber string `json:"phoneNumber" db:"phone_number"`
	Base
}

// Below are the structures of the request/response structs in the users handler.

type AddUserRequest struct {
	FirstName   string `json:"firstName" validate:"required"`
	LastName    string `json:"lastName" validate:"required"`
	Email       string `json:"email" validate:"required"`
	PhoneNumber string `json:"phoneNumber" validate:"required"`
}

type GetUserResponse struct {
	ID          int    `json:"id"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
}

type UpdateUserRequest struct {
	ID          int    `json:"id" validate:"required"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
}

type IDResponse struct {
	ID int `json:"id"`
}
