package google

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"net/http"
	"quicksend/internal/auth/jwt"
	"quicksend/internal/config"
	"quicksend/internal/subscription"
	"quicksend/internal/token"
	"quicksend/internal/user"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	googleoauth "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
)

type Source string

const (
	SourceWebsite   Source = "website"
	SourceExtension Source = "extension"
)

type Service struct {
	cfg                 *config.Config
	userService         *user.Service
	jwtService          *jwt.Service
	tokenService        *token.Service
	subscriptionService *subscription.Service
}

func NewService(
	cfg *config.Config,
	userService *user.Service,
	jwtService *jwt.Service,
	tokenService *token.Service,
) *Service {
	return &Service{
		cfg:          cfg,
		userService:  userService,
		jwtService:   jwtService,
		tokenService: tokenService,
	}
}

func (s *Service) Login(c *gin.Context) {
	source := Source(c.DefaultQuery("source", "website"))
	lang := c.DefaultQuery("lang", "en")

	oauthConfig := s.createOauthConfig(source)
	state := rand.Text()

	session := sessions.Default(c)
	session.Set("source", source)
	session.Set("lang", lang)
	session.Set("state", state)
	err := session.Save()
	if err != nil {
		log.Printf("Error saving session: %v", err)
		return
	}

	url := oauthConfig.AuthCodeURL(
		state,
		oauth2.AccessTypeOffline,
		oauth2.SetAuthURLParam("prompt", "consent"),
	)

	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (s *Service) Callback(c *gin.Context) {
	session := sessions.Default(c)

	state, _ := session.Get("state").(string)
	if state == "" || state != c.Query("state") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid oauth state"})
		return
	}

	source := Source(session.Get("source").(string))
	lang, _ := session.Get("lang").(string)

	if lang == "" {
		lang = "en"
	}

	oauthCfg := s.createOauthConfig(source)

	oauthToken, err := oauthCfg.Exchange(context.Background(), c.Query("code"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no google token data"})
		return
	}

	userInfo, err := s.getUserInfo(oauthToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "could not get user info"})
		return
	}

	u, err := s.userService.FindOrCreate(user.FindOrCreate{
		Email:      userInfo.Email,
		FirstName:  userInfo.GivenName,
		LastName:   userInfo.FamilyName,
		PictureUrl: userInfo.Picture,
		OauthID:    userInfo.Id,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if source == SourceExtension {
		_, err = s.tokenService.FindOrCreate(token.FindOrCreate{
			User:    u,
			Access:  oauthToken.AccessToken,
			Refresh: oauthToken.RefreshToken,
			Expiry:  oauthToken.Expiry,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if err := s.subscriptionService.CreateTrial(u); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	accessToken, err := s.jwtService.CreateAccessToken(u)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	refreshToken, err := s.jwtService.CreateRefreshToken(u)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	s.redirect(c, accessToken, refreshToken, source, lang)
}

func (s *Service) createOauthConfig(source Source) *oauth2.Config {
	scopes := s.cfg.WebsiteScopes
	if source == SourceExtension {
		scopes = s.cfg.ExtensionScopes
	}

	return &oauth2.Config{
		ClientID:     s.cfg.GoogleClientID,
		ClientSecret: s.cfg.GoogleClientSecret,
		RedirectURL:  s.cfg.GoogleRedirectURI(),
		Scopes:       scopes,
		Endpoint:     google.Endpoint,
	}
}

func (s *Service) getUserInfo(token *oauth2.Token) (*googleoauth.Userinfo, error) {
	httpClient := oauth2.NewClient(context.Background(),
		oauth2.StaticTokenSource(token),
	)

	svc, err := googleoauth.NewService(context.Background(),
		option.WithHTTPClient(httpClient),
	)
	if err != nil {
		return nil, err
	}

	info, err := svc.Userinfo.Get().Do()
	if err != nil {
		return nil, err
	}

	return info, nil
}

func (s *Service) redirect(
	c *gin.Context,
	accessToken, refreshToken string,
	source Source,
	lang string,
) {
	if source == SourceWebsite {
		c.SetCookie(
			"access_jwt_token",
			"Bearer "+accessToken,
			s.cfg.JWTAccessExpHours*3600,
			"/", "", true, true,
		)
		c.SetCookie(
			"refresh_jwt_token",
			"Bearer "+refreshToken,
			s.cfg.JWTRefreshExpDays*86400,
			"/", "", true, true,
		)
		c.Redirect(http.StatusTemporaryRedirect,
			fmt.Sprintf("%s/%s/profile", s.cfg.FrontendURL, lang),
		)
	} else {
		c.Redirect(http.StatusTemporaryRedirect,
			fmt.Sprintf("https://%s.chromiumapp.org/callback?access_jwt_token=%s&refresh_jwt_token=%s",
				s.cfg.ExtensionID, accessToken, refreshToken,
			),
		)
	}
}
