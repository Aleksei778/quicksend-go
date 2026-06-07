package gmail

import (
	"context"
	"fmt"
	"quicksend/internal/config"
	tokenmod "quicksend/internal/token"
	usermod "quicksend/internal/user"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

type Service struct {
	tokenSvc *tokenmod.Service
	cfg      *config.Config
}

func NewService(tokenSvc *tokenmod.Service, cfg *config.Config) *Service {
	return &Service{tokenSvc: tokenSvc, cfg: cfg}
}

func (s *Service) SendEmail(ctx context.Context, user *usermod.User, raw string) (*gmail.Message, error) {
	token, err := s.tokenSvc.FindByUser(user)
	if err != nil {
		return nil, fmt.Errorf("gmail: google token not found for user %d: %w", user.ID, err)
	}

	if token.IsExpired() {
		if err := s.tokenSvc.RefreshToken(ctx, token); err != nil {
			return nil, fmt.Errorf("gmail: google token not refreshed for user %d: %w", user.ID, err)
		}
	}

	gmailClient, err := s.getGmailClientForToken(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("gmail: google token not found for user %d: %w", user.ID, err)
	}

	return gmailClient.Users.Messages.Send("me", &gmail.Message{Raw: raw}).Context(ctx).Do()
}

func (s *Service) getGmailClientForToken(ctx context.Context, token *tokenmod.Token) (*gmail.Service, error) {
	credentials := s.createCredentials(token, []string{"https://www.googleapis.com/auth/gmail.send"})

	svc, err := gmail.NewService(ctx, option.WithTokenSource(credentials))
	if err != nil {
		return nil, fmt.Errorf("failed to create gmail service: %w", err)
	}

	return svc, nil
}

func (s *Service) createCredentials(token *tokenmod.Token, scopes []string) oauth2.TokenSource {
	oauthToken := &oauth2.Token{
		AccessToken:  token.Access,
		RefreshToken: token.Refresh,
		Expiry:       token.Expiry,
	}

	oauthCfg := oauth2.Config{
		ClientID:     s.cfg.GoogleClientID,
		ClientSecret: s.cfg.GoogleClientSecret,
		Scopes:       scopes,
		Endpoint:     google.Endpoint,
	}

	return oauthCfg.TokenSource(context.Background(), oauthToken)
}
