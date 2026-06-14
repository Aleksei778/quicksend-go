package sheets

import (
	"context"
	"fmt"
	"quicksend/internal/config"
	"quicksend/internal/google/creds"
	tokenmod "quicksend/internal/token"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type Service struct {
	tokenSvc *tokenmod.Service
	cfg      *config.Config
}

func (s *Service) getSheetsClientForToken(ctx context.Context, token *tokenmod.Token) (*sheets.Service, error) {
	credentials := creds.CreateCredentials(s.cfg, token, []string{"https://www.googleapis.com/auth/spreadsheets.readonly"})

	sheetsClient, err := sheets.NewService(ctx, option.WithTokenSource(credentials))
	if err != nil {
		return nil, fmt.Errorf("failed to create sheets service: %w", err)
	}

	return sheetsClient, nil
}
