package repository

import (
	"database/sql"
	"time"
)

type RefreshSession struct {
	Email        string
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

func (r *AuthRepository) DeleteRefreshSession(email string, fingerprint string) error {
	_, err := r.db.Exec(
		`WITH 
			expiredRows AS (
				SELECT id 
				FROM refreshSessions 
				WHERE exp <= $1
			),
			thisFingerPrint AS (
				SELECT id
				FROM refreshSessions
				WHERE user_email = $2
					AND fingerprint = $3 
			)
		DELETE FROM refreshSessions 
		WHERE id IN (
				SELECT id FROM expiredRows
				UNION ALL
				SELECT id FROM thisFingerPrint
			)`,
		time.Now(), email, fingerprint)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(`
		WITH 
			allNumeredSessions AS (
				SELECT ROW_NUMBER() over(ORDER BY exp DESC) as number, id
		 		FROM refreshSessions
				WHERE user_email = $1
			)
		DELETE FROM refreshSessions
		WHERE id IN (SELECT id FROM allNumeredSessions WHERE number > 4)
	`, email)
	return err
}

func (r *AuthRepository) WriteRefreshSession(s *RefreshSession) error {
	_, err := r.db.Exec("INSERT INTO refreshSessions(user_email, fingerprint, refreshtoken, exp, iat) VALUES($1, $2, $3, $4, $5)", s.Email, s.Fingerprint, s.RefreshToken, s.Exp, s.Iat)
	return err
}
