package gmail

import (
	"context"
	"fmt"
	"quicksend/internal/config"
	"quicksend/internal/google/creds"
	tokenmod "quicksend/internal/token"
	usermod "quicksend/internal/user"

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

func (svc *Service) SendEmail(ctx context.Context, user *usermod.User, raw string) (*gmail.Message, error) {
	token, err := svc.tokenSvc.FindByUser(user)
	if err != nil {
		return nil, fmt.Errorf("gmail: google token not found for user %d: %w", user.ID, err)
	}

	if token.IsExpired() {
		if err := svc.tokenSvc.RefreshToken(ctx, token); err != nil {
			return nil, fmt.Errorf("gmail: google token not refreshed for user %d: %w", user.ID, err)
		}
	}

	gmailClient, err := svc.getGmailClientForToken(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("gmail: google token not found for user %d: %w", user.ID, err)
	}

	return gmailClient.Users.Messages.Send("me", &gmail.Message{Raw: raw}).Context(ctx).Do()
}

func (svc *Service) getGmailClientForToken(ctx context.Context, token *tokenmod.Token) (*gmail.Service, error) {
	credentials := creds.CreateCredentials(svc.cfg, token, []string{"https://www.googleapis.com/auth/gmail.send"})

	gmailClient, err := gmail.NewService(ctx, option.WithTokenSource(credentials))
	if err != nil {
		return nil, fmt.Errorf("failed to create gmail service: %w", err)
	}

	return gmailClient, nil
}
