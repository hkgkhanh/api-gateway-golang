package storage

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"api-gateway-golang/internal/model"
)

func (s *Storage) AddProduct(ctx context.Context, product model.AddProductRequest) (int, error) {
	var id int
	err := s.db.Get(&id, `INSERT INTO products(name, manufacturer, quantity, created_at, updated_at)
			VALUES($1,$2,$3,$4,$5) RETURNING id`, product.Name, product.Manufacturer, product.Quantity, time.Now().UTC(), time.Now().UTC())

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Storage) GetProduct(ctx context.Context, id int) (model.Product, error) {
	var product model.Product

	err := s.db.Get(&product, `Select * from products where id=$1`, id)
	if err != nil {
		return product, err
	}

	return product, nil
}

func (s *Storage) GetProducts(ctx context.Context) ([]model.Product, error) {
	var products []model.Product
	err := s.db.Select(&products, `SELECT * from products`)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (s *Storage) UpdateProduct(ctx context.Context, product model.UpdateProductRequest) (int, error) {
	var columns []string
	var argCount = 1
	var args []interface{}

	if product.Name != "" {
		columns = append(columns, fmt.Sprintf("name=$%d", argCount))
		args = append(args, product.Name)
		argCount++
	}

	if product.Manufacturer != "" {
		columns = append(columns, fmt.Sprintf("manufacturer=$%d", argCount))
		args = append(args, product.Manufacturer)
		argCount++
	}

	if product.Quantity < 0 {
		columns = append(columns, fmt.Sprintf("quantity=$%d", argCount))
		args = append(args, product.Quantity)
		argCount++
	}

	columns = append(columns, fmt.Sprintf("updated_at=$%d", argCount))
	args = append(args, time.Now().UTC())
	argCount++

	if len(columns) == 0 {
		return 0, errors.New("no fields to update")
	}

	args = append(args, product.ID)

	query := fmt.Sprintf(`UPDATE products SET %s WHERE id=$%d RETURNING id`, strings.Join(columns, ", "), argCount)

	var id int
	err := s.db.Get(&id, query, args...)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *Storage) DeleteProduct(ctx context.Context, id int) error {
	_, err := s.db.Exec(`DELETE FROM products WHERE id=$1`, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) VerifyProductExists(ctx context.Context, id int) (bool, error) {
	var exists bool
	err := s.db.Get(&exists, `SELECT EXISTS(SELECT 1 from products where id=$1)`, id)
	if err != nil {
		return false, err
	}

	return exists, nil
}
