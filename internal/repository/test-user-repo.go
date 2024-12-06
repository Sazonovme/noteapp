package repository

import (
	"database/sql"
	"noteapp/internal/model"
)

type TestUserRepository struct {
	db *sql.DB
}

func NewTestUserRepository(db *sql.DB) *TestUserRepository {
	return &TestUserRepository{
		db: db,
	}
}

func (r *TestUserRepository) CreateUser(u *model.User) error {
	_, err := r.db.Exec(
		"INSERT INTO test_users(login, password) VALUES ($1, $2)",
		u.Login, u.Password,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *TestUserRepository) FindByLogin(login string) (*model.User, error) {
	u := model.User{}
	if err := r.db.QueryRow(
		"SELECT id, login, password FROM test_users WHERE login=$1",
		login,
	).Scan(
		&u.Id,
		&u.Login,
		&u.Password,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &u, nil
}
