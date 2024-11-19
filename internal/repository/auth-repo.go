package repository

import (
	"database/sql"
	"time"
)

type RefreshSession struct {
	Login        string
	RefreshToken string
	Exp          time.Time
	Iat          time.Time
	Fingerprint  string
}

type AuthRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{
		db: db,
	}
}

func (r *AuthRepository) DeleteRefreshSession(login string, fingerprint string) error {
	_, err := r.db.Exec("DELETE FROM refreshSessions WHERE user_login = $1 AND fingerprint = $2", login, fingerprint)
	return err
}

func (r *AuthRepository) WriteRefreshSession(s *RefreshSession) error {
	_, err := r.db.Exec("INSERT INTO refreshSessions VALUES($1, $2, $3, $4, $5)", s.Login, s.RefreshToken, s.Exp, s.Iat, s.Fingerprint)
	return err
}

func (r *AuthRepository) RefreshSessionExist(rToken string, fingerprint string) error {
	_, err := r.db.Query("SELECT login, fingerprint FROM refreshSessions WHERE refreshtoken = $1 AND fingerprint = $2", rToken, fingerprint)
	return err
}
