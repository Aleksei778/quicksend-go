package google

import (
	"crypto/rand"
	"log"
	"net/http"
	"quicksend/internal/auth/jwt"
	"quicksend/internal/config"
	"quicksend/internal/token"
	"quicksend/internal/user"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Source string

const (
	SourceWebsite   Source = "website"
	SourceExtension Source = "extension"
)

type Service struct {
	cfg          *config.Config
	userService  *user.Service
	jwtService   *jwt.Service
	tokenService *token.Service
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

	oauthToken, err := oauthCfg.Exchange(c.Background(), c.Query("code"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no google token data"})
		return
	}

	userInfo, err := s.getUser

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
