package sqlstore

import (
	"database/sql"
	"errors"
	"noteapp/internal/noteapp/user"
)

type UserRepository struct {
	store *Sqlstore
}

func NewUserRepository(store *Sqlstore) *UserRepository {
	return &UserRepository{
		store: store,
	}
}

////////// SQL Requests for User //////////

func (r *UserRepository) CreateUser(u *user.User) error {
	_, err := r.FindByLogin(u.Login)
	if err != nil {

		// encrypt password
		if _, err = u.EncryptPassword(); err != nil {
			return err
		}

		_, err := r.store.Db.Exec(
			"INSERT INTO users(login, password) VALUES ($1, $2)",
			u.Login, u.Password,
		)
		if err != nil {
			return err
		}

		// Создать рефрешь токен
		return nil
	}
	return errors.New("user alredy exist")
}

func (r *UserRepository) FindByLogin(login string) (*user.User, error) {

	u := user.New()
	if err := r.store.Db.QueryRow(
		"SELECT id, login, password FROM users WHERE login=$1",
		login,
	).Scan(
		&u.Id,
		&u.Login,
		&u.Password,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return u, nil
}

// func (r *UserRepository) AuthenticateUser(login string) error {

// 	// На 24 часа вперед
// 	t := time.Now()
// 	t = t.Add(24 * time.Hour)

// 	u := user.New()
// 	err := r.store.Db.QueryRow(
// 		"SELECT login FROM refreshTokens WHERE login = $1 AND exp < $2", login,
// 	).Scan(
// 		&u.Login,
// 	)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			// создаем

// 		}
// 		return err
// 	}
// 	// Создать новую пару токенов и записать в базу
// 	// Вернуть клиенту
// }
