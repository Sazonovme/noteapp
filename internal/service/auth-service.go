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
	ErrTokenInvalid         = errors.New("invalid token or token time has expired")
)

type RequestTokenData struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
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

func (s *AuthService) MakeRefreshSession(email string, fingerprint string) (*RequestTokenData, error) {

	err := s.repository.DeleteRefreshSession(email, fingerprint)
	if err != nil {
		logger.NewLog("service - MakeRefreshSession()", 2, err, "Filed to delete session in repository", "email: "+email+", fingerprint: "+fingerprint)
		return nil, err
	}

	expAccess := time.Now().Add(15 * time.Minute)
	aToken, err := CreateJWTToken(email, expAccess, secretKeyAccess)
	if err != nil {
		logger.NewLog("service - MakeRefreshSession()", 2, err, "Filed to Create JWT Token (Access)", "login: "+email+", expAccess: "+expAccess.Format("2006-01-02 15:04"))
		return nil, err
	}

	expRefresh := time.Now().Add(15 * 24 * time.Hour)
	rToken, err := CreateJWTToken(email, expRefresh, secretKeyRefresh)
	if err != nil {
		logger.NewLog("service - MakeRefreshSession()", 2, err, "Filed to Create JWT Token (Refresh)", "login: "+email+", expRefresh: "+expRefresh.Format("2006-01-02 15:04"))
		return nil, err
	}

	session := repository.RefreshSession{
		Email:        email,
		RefreshToken: rToken,
		Exp:          expRefresh,
		Iat:          time.Now(),
		Fingerprint:  fingerprint,
	}
	err = s.repository.WriteRefreshSession(&session)
	if err != nil {
		logger.NewLog("service - MakeRefreshSession()", 2, err, "Filed to write refresh session in repository", session)
		return nil, err
	}

	return &RequestTokenData{
		AccessToken:  aToken,
		RefreshToken: rToken,
	}, nil
}

func (s *AuthService) UpdateTokens(oldRefreshToken string, fingerprint string) (*RequestTokenData, error) {
	email, err := VerifyToken(oldRefreshToken, secretKeyRefresh)
	if err != nil {
		logger.NewLog("service - UpdateTokens()", 5, err, "Filed to verify token", nil)
		return nil, err
	}

	err = s.repository.DeleteRefreshSession(email, fingerprint)
	if err != nil {
		logger.NewLog("service - UpdateTokens()", 2, err, "Filed delete refresh session in repository", "email: "+email+", fingerprint: "+fingerprint)
		return nil, err
	}

	refreshSesssion, err := s.MakeRefreshSession(email, fingerprint)
	if err != nil {
		logger.NewLog("service - UpdateTokens()", 2, err, "Filed to make refresh session", "email: "+email+", fingerprint: "+fingerprint)
		return nil, err
	}

	return refreshSesssion, nil
}

func CreateJWTToken(email string, exp time.Time, secretKey string) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ExpiresAt: &jwt.NumericDate{Time: exp},
		Subject:   email,
		IssuedAt:  &jwt.NumericDate{Time: time.Now()},
	})
	tokenString, err := claims.SignedString([]byte(secretKey))
	if err != nil {
		logger.NewLog("service - CreateJWTToken()", 2, err, "Filed to sign JWT token", nil)
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(tokenString string, secretKey string) (string, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		logger.NewLog("service - VerifyToken()", 2, err, "Filed to parse JWT token", nil)
		if err.Error() == "Token is expired" {
			return "", ErrTokenInvalid
		}
		return "", err
	}

	if !token.Valid {
		return "", ErrTokenInvalid
	}

	email, ok := claims["sub"].(string)
	if !ok {
		logger.NewLog("service - VerifyToken()", 2, err, "Field sub(login) not exist in token", nil)
		return "", ErrTokenInvalid
	}

	return email, nil
}

func VerifyAccessToken(tokenString string) (string, error) {
	return VerifyToken(tokenString, secretKeyAccess)
}
