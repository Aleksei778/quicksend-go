package sheets

import (
	"context"
	"fmt"
	"quicksend/internal/config"
	"quicksend/internal/google/creds"
	tokenmod "quicksend/internal/token"
	"strings"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type Service struct {
	tokenSvc *tokenmod.Service
	cfg      *config.Config
}

func NewService(tokenSvc *tokenmod.Service, cfg *config.Config) *Service {
	return &Service{tokenSvc: tokenSvc, cfg: cfg}
}

func (svc *Service) ParseEmails(ctx context.Context, token *tokenmod.Token, spreadsheetID, rang string) ([]string, error) {
	sheetsClient, err := svc.getSheetsClientForToken(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("could not get sheets client: %w", err)
	}

	result, err := sheetsClient.Spreadsheets.Values.Get(spreadsheetID, rang).Do()
	if err != nil {
		return nil, fmt.Errorf("could not retrieve data from sheets: %w", err)
	}

	if len(result.Values) == 0 {
		return nil, fmt.Errorf("could not retrieve data from sheets: %w", err)
	}

	emails := svc.extractEmails(result.Values)
	if len(emails) == 0 {
		return nil, fmt.Errorf("could not retrieve data from sheets: %w", err)
	}

	return emails, nil
}

func (svc *Service) getSheetsClientForToken(ctx context.Context, token *tokenmod.Token) (*sheets.Service, error) {
	credentials := creds.CreateCredentials(svc.cfg, token, []string{"https://www.googleapis.com/auth/spreadsheets.readonly"})

	sheetsClient, err := sheets.NewService(ctx, option.WithTokenSource(credentials))
	if err != nil {
		return nil, fmt.Errorf("failed to create sheets service: %w", err)
	}

	return sheetsClient, nil
}

func (svc *Service) extractEmails(values [][]interface{}) []string {
	seen := make(map[string]struct{})
	emails := make([]string, 0)

	for _, row := range values {
		if len(row) == 0 {
			continue
		}

		email, ok := row[0].(string)
		if !ok || !strings.Contains(email, "@") {
			continue
		}

		if _, exists := seen[email]; !exists {
			seen[email] = struct{}{}
			emails = append(emails, email)
		}
	}

	return emails
}
