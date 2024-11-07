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
	exp          string
}

type UserClaims struct {
	Subr      string `json:"subr"`
	RegClaims jwt.RegisteredClaims
}

func MakeRefreshSession(db *sql.DB, login string, fingerprint string) (*RefreshSession, error) {

	// Удалить старые токены если такие были
	db.Exec("DELETE FROM refreshTokens WHERE login = $1 AND fingerprint = $2", login, fingerprint)

	aToken, expAccess, err := createAccessToken(login)
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
		exp:          expAccess.Time.Format("2006-01-02 15:04"),
	}, nil
}

func UpdateTokens(db *sql.DB, oldRefreshToken string, oldFingerPrint string) (*RefreshSession, error) {

	// Проверка валидности токена
	claims, err := verifyToken(oldRefreshToken, secretKeyRefresh)
	if err != nil {
		return nil, err
	}

	expTime := checkExpiredTime(claims)
	if !expTime {
		return nil, errors.New("token lifetime expired")
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

func createAccessToken(login string) (string, *jwt.NumericDate, error) {

	lifeTime := jwt.NewNumericDate(time.Now().Add(15 * time.Minute))

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": login,             // Subject (user identifier)
		"exp": lifeTime,          // Expiration time
		"iat": time.Now().Unix(), // Время выпуска
	})

	// Подписание токена
	tokenString, err := claims.SignedString([]byte(secretKeyAccess))
	if err != nil {
		return "", nil, err
	}

	return tokenString, lifeTime, nil
}

func verifyToken(tokenString string, secretKey string) (jwt.MapClaims, error) {

	// Parse the token with the secret key
	claims := jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
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
	return claims, nil
}

func VerifyAccessToken(tokenString string) error {
	claims, err := verifyToken(tokenString, secretKeyAccess)
	if err != nil {
		return err
	}

	expTime := checkExpiredTime(claims)

	if !expTime {
		return errors.New("token lifetime has expired")
	}
	return nil
}

func checkExpiredTime(claims jwt.MapClaims) bool {
	var tm time.Time
	switch exp := claims["exp"].(type) {
	case float64:
		tm = time.Unix(int64(exp), 0)
	case json.Number:
		v, _ := exp.Int64()
		tm = time.Unix(v, 0)
	}
	//fmt.Println("exp:", tm.Add(-30*time.Second), "now", time.Now())
	return tm.Add(-30*time.Second).Unix() >= time.Now().Unix()
}
