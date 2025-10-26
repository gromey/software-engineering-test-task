package repository

import (
	"context"
	"database/sql"
	"errors"

	"cruder/internal/model"
	"cruder/pkg/validation"
)

type UserRepository interface {
	GetAll(ctx context.Context) ([]model.User, error)
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	GetByID(ctx context.Context, id int64) (*model.User, error)
	Post(ctx context.Context, user *model.User) (int64, error)
	Patch(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id int64) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

const getAllStm = `SELECT id, username, email, full_name FROM users`

func (r *userRepository) GetAll(ctx context.Context) ([]model.User, error) {
	rows, err := r.db.QueryContext(ctx, getAllStm)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) { _ = rows.Close() }(rows)

	var users []model.User
	for rows.Next() {
		var u model.User
		if err = rows.Scan(&u.ID, &u.Username, &u.Email, &u.FullName); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

const getByUsernameStm = `SELECT id, username, email, full_name FROM users WHERE username = $1`

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	var u model.User
	if err := r.db.QueryRowContext(ctx, getByUsernameStm, username).
		Scan(&u.ID, &u.Username, &u.Email, &u.FullName); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, validation.ErrUserNotFound
		}
		return nil, err
	}
	return &u, nil
}

const getByIDStm = `SELECT id, username, email, full_name FROM users WHERE id = $1`

func (r *userRepository) GetByID(ctx context.Context, id int64) (*model.User, error) {
	var u model.User
	if err := r.db.QueryRowContext(ctx, getByIDStm, id).
		Scan(&u.ID, &u.Username, &u.Email, &u.FullName); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, validation.ErrUserNotFound
		}
		return nil, err
	}
	return &u, nil
}

const postStm = `INSERT INTO users (username, email, full_name) VALUES ($1, $2, $3) RETURNING id`

func (r *userRepository) Post(ctx context.Context, user *model.User) (int64, error) {
	var id int64
	if err := r.db.QueryRowContext(ctx, postStm, user.Username, user.Email, user.FullName).
		Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

const patchStm = `UPDATE users SET username = $1, email = $2, full_name = $3 WHERE id = $4`

func (r *userRepository) Patch(ctx context.Context, user *model.User) error {
	_, err := r.db.ExecContext(ctx, patchStm, user.Username, user.Email, user.FullName, user.ID)
	return err
}

const deleteStm = `DELETE FROM users WHERE id = $1`

func (r *userRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, deleteStm, id)
	return err
}
