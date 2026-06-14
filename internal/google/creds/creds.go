package creds

import (
	"context"
	"quicksend/internal/config"
	tokenmod "quicksend/internal/token"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func CreateCredentials(cfg *config.Config, token *tokenmod.Token, scopes []string) oauth2.TokenSource {
	oauthToken := &oauth2.Token{
		AccessToken:  token.Access,
		RefreshToken: token.Refresh,
		Expiry:       token.Expiry,
	}

	oauthCfg := oauth2.Config{
		ClientID:     cfg.GoogleClientID,
		ClientSecret: cfg.GoogleClientSecret,
		Scopes:       scopes,
		Endpoint:     google.Endpoint,
	}

	return oauthCfg.TokenSource(context.Background(), oauthToken)
}
