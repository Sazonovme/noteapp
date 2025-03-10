package repository

import (
	"database/sql"
	"errors"
	"noteapp/internal/model"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) CreateUser(u *model.User) error {
	_, err := r.db.Exec(
		"INSERT INTO users(email, password) VALUES ($1, $2)",
		u.Email, u.Password,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) FindByLogin(email string) (*model.User, error) {
	u := model.User{}
	if err := r.db.QueryRow(
		"SELECT id, email, password FROM users WHERE email=$1",
		email,
	).Scan(
		&u.Id,
		&u.Email,
		&u.Password,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &u, nil
}
