package storage

import (
	"api-gateway-golang/internal/model"
	"context"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Initializes the postgres driver
)

type StorageInterface interface {
	AddUser(ctx context.Context, user model.AddUserRequest) (int, error)
	GetUser(ctx context.Context, id int) (model.User, error)
	GetUsers(ctx context.Context) ([]model.User, error)
	UpdateUser(ctx context.Context, user model.UpdateUserRequest) (int, error)
	DeleteUser(ctx context.Context, id int) error
	VerifyUserExists(ctx context.Context, id int) (bool, error)

	AddProduct(ctx context.Context, product model.AddProductRequest) (int, error)
	GetProduct(ctx context.Context, id int) (model.Product, error)
	GetProducts(ctx context.Context) ([]model.Product, error)
	UpdateProduct(ctx context.Context, product model.UpdateProductRequest) (int, error)
	DeleteProduct(ctx context.Context, id int) error
	VerifyProductExists(ctx context.Context, id int) (bool, error)
}

// Storage contains an SQL db. Storage implements the StorageInterface.
type Storage struct {
	db *sqlx.DB
}

func (s *Storage) Close() error {
	if err := s.db.Close(); err != nil {
		return err
	}

	return nil
}

func (s *Storage) GetDB() *sqlx.DB {
	return s.db
}
