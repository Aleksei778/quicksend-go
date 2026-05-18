package jwt

import (
	"errors"
	"quicksend/internal/config"
	"quicksend/internal/user"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID  uint   `json:"user_id"`
	Email   string `json:"email"`
	OauthID string `json:"oauth_id"`
	Type    string `json:"type"`
	jwt.RegisteredClaims
}

type Service struct {
	cfg *config.Config
}

func NewService(cfg *config.Config) *Service {
	return &Service{cfg: cfg}
}

func (s *Service) CreateAccessToken(u *user.User) (string, error) {
	return s.sign(u, "access", time.Duration(s.cfg.JWTAccessExpHours)*time.Hour, s.cfg.JWTAccessSecret)
}

func (s *Service) CreateRefreshToken(u *user.User) (string, error) {
	return s.sign(u, "refresh", time.Duration(s.cfg.JWTRefreshExpDays)*24*time.Hour, s.cfg.JWTRefreshSecret)
}

func (s *Service) VerifyAccessToken(tokenStr string) (*Claims, error) {
	return s.verify(tokenStr, "access", s.cfg.JWTAccessSecret)
}

func (s *Service) VerifyRefreshToken(tokenStr string) (*Claims, error) {
	return s.verify(tokenStr, "refresh", s.cfg.JWTRefreshSecret)
}

func (s *Service) sign(u *user.User, tokenType string, ttl time.Duration, secret string) (string, error) {
	claims := Claims{
		UserID:  u.ID,
		Email:   u.Email,
		OauthID: u.OauthID,
		Type:    tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
}

func (s *Service) verify(tokenStr string, expectedType string, secret string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}

		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	if claims.Type != expectedType {
		return nil, errors.New("invalid token type")
	}

	return claims, nil
}
