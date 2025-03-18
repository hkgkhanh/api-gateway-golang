package storage

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"api-gateway-golang/internal/model"
)

func (s *Storage) AddUser(ctx context.Context, user model.AddUserRequest) (int, error) {
	var id int
	err := s.db.Get(&id, `INSERT INTO users(first_name, last_name, email, phone_number, created_at, updated_at)
			VALUES($1,$2,$3,$4,$5,$6) RETURNING id`, user.FirstName, user.LastName, user.Email, user.PhoneNumber, time.Now().UTC(), time.Now().UTC())

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Storage) GetUser(ctx context.Context, id int) (model.User, error) {
	var user model.User

	err := s.db.Get(&user, `Select * from users where id=$1`, id)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (s *Storage) GetUsers(ctx context.Context) ([]model.User, error) {
	var users []model.User
	err := s.db.Select(&users, `SELECT * from users`)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *Storage) UpdateUser(ctx context.Context, user model.UpdateUserRequest) (int, error) {
	var columns []string
	var argCount = 1
	var args []interface{}

	if user.FirstName != "" {
		columns = append(columns, fmt.Sprintf("first_name=$%d", argCount))
		args = append(args, user.FirstName)
		argCount++
	}

	if user.LastName != "" {
		columns = append(columns, fmt.Sprintf("last_name=$%d", argCount))
		args = append(args, user.LastName)
		argCount++
	}

	if user.Email != "" {
		columns = append(columns, fmt.Sprintf("email=$%d", argCount))
		args = append(args, user.Email)
		argCount++
	}

	if user.PhoneNumber != "" {
		columns = append(columns, fmt.Sprintf("phone_number=$%d", argCount))
		args = append(args, user.PhoneNumber)
		argCount++
	}

	columns = append(columns, fmt.Sprintf("updated_at=$%d", argCount))
	args = append(args, time.Now().UTC())
	argCount++

	if len(columns) == 0 {
		return 0, errors.New("no fields to update")
	}

	args = append(args, user.ID)

	query := fmt.Sprintf(`UPDATE users SET %s WHERE id=$%d RETURNING id`, strings.Join(columns, ", "), argCount)

	var id int
	err := s.db.Get(&id, query, args...)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *Storage) DeleteUser(ctx context.Context, id int) error {
	_, err := s.db.Exec(`DELETE FROM users WHERE id=$1`, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) VerifyUserExists(ctx context.Context, id int) (bool, error) {
	var exists bool
	err := s.db.Get(&exists, `SELECT EXISTS(SELECT 1 from users where id=$1)`, id)
	if err != nil {
		return false, err
	}

	return exists, nil
}
