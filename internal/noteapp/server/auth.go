package server

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	secretKeyRefresh string = "super-secret-key-refresh"
	secretKeyAccess  string = "super-secret-key-access"
)

type RefreshSession struct {
	accessToken  string
	refreshToken string
}

func MakeRefreshSession(db *sql.DB, login string, fingerprint string) (*RefreshSession, error) {
	aToken, err := createAccessToken(login)
	if err != nil {
		return nil, err
	}
	rToken, exp, iat, err := createRefreshToken(login)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("INSERT INTO refreshTokens VALUES($1, $2, $3, $4, $5)", login, rToken, exp, iat, fingerprint)
	if err != nil {
		return nil, err
	}

	return &RefreshSession{
		accessToken:  aToken,
		refreshToken: rToken,
	}, nil
}

func UpdateTokens(db *sql.DB, oldRefreshToken string, oldFingerPrint string) (*RefreshSession, error) {

	// Проверка валидности токена
	oldJWT, err := verifyToken(oldRefreshToken, secretKeyRefresh)
	if err != nil {
		return nil, err
	}

	// Проверка времени токена
	claims, ok := oldJWT.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("can't convert token's claims to standard claims")
	}

	var tm time.Time
	switch exp := claims["exp"].(type) {
	case float64:
		tm = time.Unix(int64(exp), 0)
	case json.Number:
		v, _ := exp.Int64()
		tm = time.Unix(v, 0)
	}

	if tm.Add(-30*time.Second).Unix() <= time.Now().Unix() {
		return nil, errors.New("token lifetime has expired")
	}

	rows, err := db.Query("DELETE FROM refreshTokens WHERE refreshtoken = $1 AND fingerprint = $2 RETURNING login, fingerprint", oldRefreshToken, oldFingerPrint)
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

	refreshSesssion, err := MakeRefreshSession(db, login, fingerprint)
	if err != nil {
		return nil, err
	}
	return refreshSesssion, nil
}

func createRefreshToken(login string) (string, string, string, error) {

	lifeTime := time.Now().Add(30 * 24 * time.Hour).Format("2006-01-02 15:04:05")
	iat := time.Now().Format("2006-01-02 15:04:05")

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": login,    // Subject (user identifier)
		"exp": lifeTime, // Expiration time
		"iat": iat,      // Время выпуска
	})

	// Подписание токена
	tokenString, err := claims.SignedString([]byte(secretKeyRefresh))
	if err != nil {
		return "", "", "", err
	}

	return tokenString, lifeTime, iat, nil
}

func createAccessToken(login string) (string, error) {

	lifeTime := time.Now().Add(15 * time.Minute).Format("2006-01-02 15:04:05")

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": login,                                    // Subject (user identifier)
		"exp": lifeTime,                                 // Expiration time
		"iat": time.Now().Format("2006-01-02 15:04:05"), // Время выпуска
	})

	// Подписание токена
	tokenString, err := claims.SignedString([]byte(secretKeyAccess))
	if err != nil {
		return "", err
	}

	return tokenString, nil

}

// Function to verify JWT tokens
func verifyToken(tokenString string, secretKey string) (*jwt.Token, error) {

	// Parse the token with the secret key

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	// Check for verification errors
	if err != nil {
		return nil, err
	}

	// Check if the token is valid
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Return the verified token
	return token, nil
}

func VerifyAccessToken(tokenString string) error {
	_, err := verifyToken(tokenString, secretKeyAccess)
	if err != nil {
		return err
	}
	return nil
}
