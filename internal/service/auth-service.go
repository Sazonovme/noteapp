package service

import (
	"errors"
	"noteapp/internal/repository"
	"noteapp/pkg/logger"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var (
	secretKeyRefresh string = "super-secret-key-refresh"
	secretKeyAccess  string = "super-secret-key-access"
	errTokenInvalid         = errors.New("invalid token or token time has expired")
)

type RequestTokenData struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	Exp          string `json:"exp"`
}

type AuthRepository interface {
	DeleteRefreshSession(string, string) error
	WriteRefreshSession(*repository.RefreshSession) error
}

type AuthService struct {
	repository AuthRepository
}

func NewAuthService(r AuthRepository) *AuthService {
	return &AuthService{
		repository: r,
	}
}

func (s *AuthService) MakeRefreshSession(login string, fingerprint string) (*RequestTokenData, error) {
	err := s.repository.DeleteRefreshSession(login, fingerprint)
	if err != nil {
		return nil, err
	}

	expAccess := time.Now().Add(15 * time.Minute)
	aToken, err := CreateJWTToken(login, expAccess, secretKeyAccess)
	if err != nil {
		return nil, err
	}

	expRefresh := time.Now().Add(15 * 24 * time.Hour)
	rToken, err := CreateJWTToken(login, expRefresh, secretKeyRefresh)
	if err != nil {
		return nil, err
	}

	session := repository.RefreshSession{
		Login:        login,
		RefreshToken: rToken,
		Exp:          expRefresh,
		Iat:          time.Now(),
		Fingerprint:  fingerprint,
	}
	err = s.repository.WriteRefreshSession(&session)
	if err != nil {
		return nil, err
	}

	return &RequestTokenData{
		AccessToken:  aToken,
		RefreshToken: rToken,
		Exp:          expAccess.Format("2006-01-02 15:04"),
	}, nil
}

func (s *AuthService) UpdateTokens(oldRefreshToken string, fingerprint string, login string) (*RequestTokenData, error) {
	_, err := VerifyToken(oldRefreshToken, secretKeyRefresh)
	if err != nil {
		return nil, err
	}

	err = s.repository.DeleteRefreshSession(login, fingerprint)
	if err != nil {
		return nil, err
	}

	refreshSesssion, err := s.MakeRefreshSession(login, fingerprint)
	if err != nil {
		return nil, err
	}

	return refreshSesssion, nil
}

func CreateJWTToken(login string, exp time.Time, secretKey string) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ExpiresAt: &jwt.NumericDate{Time: exp},
		Subject:   login,
		IssuedAt:  &jwt.NumericDate{Time: time.Now()},
	})
	tokenString, err := claims.SignedString([]byte(secretKey))
	if err != nil {
		logger.NewLog("service - CreateJWTToken()", 2, err, "Filed to sign JWT token", nil)
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(tokenString string, secretKey string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		logger.NewLog("service - VerifyToken()", 2, err, "Filed to parse JWT token", nil)
		return false, err
	}

	if !token.Valid {
		return false, errTokenInvalid
	}

	return true, nil
}

func VerifyAccessToken(tokenString string) (bool, error) {
	return VerifyToken(tokenString, secretKeyAccess)
}
