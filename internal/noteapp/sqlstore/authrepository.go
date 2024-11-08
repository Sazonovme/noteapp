package sqlstore

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var (
	secretKeyRefresh string = "super-secret-key-refresh"
	secretKeyAccess  string = "super-secret-key-access"
)

type AuthRepository struct {
	store *Sqlstore
}

type RefreshSession struct {
	AccessToken  string
	RefreshToken string
	Exp          string
}

func NewAuthRepository(store *Sqlstore) *AuthRepository {
	return &AuthRepository{
		store: store,
	}
}

func (r *AuthRepository) MakeRefreshSession(login string, fingerprint string) (*RefreshSession, error) {

	r.store.Db.Exec("DELETE FROM refreshTokens WHERE login = $1 AND fingerprint = $2", login, fingerprint)

	expAccess := time.Now().Add(15 * time.Minute)
	aToken, err := createJWTToken(login, expAccess, secretKeyAccess)
	if err != nil {
		return nil, err
	}

	expRefresh := time.Now().Add(15 * time.Minute)
	rToken, err := createJWTToken(login, expRefresh, secretKeyRefresh)
	if err != nil {
		return nil, err
	}

	_, err = r.store.Db.Exec("INSERT INTO refreshTokens VALUES($1, $2, $3, $4, $5)", login, rToken, expRefresh, time.Now(), fingerprint)
	if err != nil {
		return nil, err
	}

	return &RefreshSession{
		AccessToken:  aToken,
		RefreshToken: rToken,
		Exp:          expAccess.Format("2006-01-02 15:04"),
	}, nil
}

func (r *AuthRepository) UpdateTokens(oldRefreshToken string, oldFingerPrint string) (*RefreshSession, error) {

	_, err := verifyToken(oldRefreshToken, secretKeyRefresh)
	if err != nil {
		return nil, err
	}

	rows, err := r.store.Db.Query("SELECT login, fingerprint FROM refreshTokens WHERE refreshtoken = $1 AND fingerprint = $2", oldRefreshToken, oldFingerPrint)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var login, fingerprint string
	if !rows.Next() {
		return nil, errors.New("token with this refreshJWT and fingerptint not exist in databse")
	} else {
		err = rows.Scan(&login, &fingerprint)
		if err != nil {
			return nil, errors.New("not successful scan login and fingerprint after DELETE SQL")
		}
	}

	refreshSesssion, err := r.MakeRefreshSession(login, fingerprint)
	if err != nil {
		return nil, err
	}
	return refreshSesssion, nil
}

func (r *AuthRepository) VerifyAccessToken(token string) error {
	_, err := verifyToken(token, secretKeyAccess)
	if err != nil {
		return err
	}
	return nil
}

func createJWTToken(login string, exp time.Time, secretKey string) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ExpiresAt: &jwt.NumericDate{Time: exp},
		Subject:   login,
		IssuedAt:  &jwt.NumericDate{Time: time.Now()},
	})
	tokenString, err := claims.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func verifyToken(tokenString string, secretKey string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return false, err
	}
	if !token.Valid {
		return false, errors.New("invalid token or token time has expired")
	}
	return true, nil
}
