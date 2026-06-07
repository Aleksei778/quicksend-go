package token

import (
	"context"
	"fmt"
	"quicksend/internal/config"
	"quicksend/internal/user"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gorm.io/gorm"
)

type Service struct {
	db   *gorm.DB
	cfg  *config.Config
	repo *Repository
}

func NewService(db *gorm.DB, repo *Repository, cfg *config.Config) *Service {
	return &Service{db: db, repo: repo, cfg: cfg}
}

func (svc *Service) RefreshToken(ctx context.Context, token *Token) error {
	if token == nil || token.Refresh == "" {
		return fmt.Errorf("token: refresh token is missing")
	}

	oauthCfg := &oauth2.Config{
		ClientID:     svc.cfg.GoogleClientID,
		ClientSecret: svc.cfg.GoogleClientSecret,
		Endpoint:     google.Endpoint,
	}

	oauthToken := &oauth2.Token{
		RefreshToken: token.Refresh,
	}

	newToken, err := oauthCfg.TokenSource(ctx, oauthToken).Token()
	if err != nil {
		return fmt.Errorf("token: %w", err)
	}

	token.Access = newToken.AccessToken
	token.Refresh = newToken.RefreshToken

	if err := svc.db.WithContext(ctx).Save(token).Error; err != nil {
		return fmt.Errorf("google_token_service:RefreshToken: failed to save token: %w", err)
	}

	return nil
}

func (svc *Service) FindOrCreate(dto FindOrCreate) (*Token, error) {
	token, err := svc.FindByUser(dto.User)

	if err != nil {
		return nil, err
	}

	if token == nil {
		return svc.create(dto)
	}

	token.Access = dto.Access
	token.Refresh = dto.Refresh
	token.Expiry = dto.Expiry

	if err := svc.db.Save(token).Error; err != nil {
		return nil, err
	}

	return token, nil
}

func (svc *Service) create(dto FindOrCreate) (*Token, error) {
	token := &Token{
		UserID:  dto.User.ID,
		Access:  dto.Access,
		Refresh: dto.Refresh,
		Expiry:  dto.Expiry,
	}

	if err := svc.db.Create(token).Error; err != nil {
		return nil, err
	}

	return token, nil
}

func (svc *Service) FindByUser(u *user.User) (*Token, error) {
	return svc.repo.FindByUser(u)
}

func (svc *Service) update(access string, expiry time.Time, token *Token) (*Token, error) {
	token.Access = access
	token.Expiry = expiry

	if err := svc.db.Save(token).Error; err != nil {
		return nil, err
	}

	return token, nil
}
