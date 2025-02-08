package repository

import "database/sql"

type TestAuthRepository struct {
	db *sql.DB
}

func NewTestAuthRepository(db *sql.DB) *TestAuthRepository {
	return &TestAuthRepository{
		db: db,
	}
}

func (r *TestAuthRepository) DeleteRefreshSession(email string, fingerprint string) error {
	_, err := r.db.Exec("DELETE FROM test_refreshSessions WHERE user_email = $1 AND fingerprint = $2", email, fingerprint)
	return err
}

func (r *TestAuthRepository) WriteRefreshSession(s *RefreshSession) error {
	_, err := r.db.Exec("INSERT INTO test_refreshSessions(user_email, fingerprint, refreshtoken, exp, iat) VALUES($1, $2, $3, $4, $5)", s.Email, s.Fingerprint, s.RefreshToken, s.Exp, s.Iat)
	return err
}

func (r *TestAuthRepository) RefreshSessionExist(rToken string, fingerprint string) error {
	_, err := r.db.Query("SELECT email, fingerprint FROM test_refreshSessions WHERE refreshtoken = $1 AND fingerprint = $2", rToken, fingerprint)
	return err
}
