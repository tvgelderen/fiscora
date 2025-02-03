package auth

import (
	"github.com/tvgelderen/fiscora/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type AuthService struct {
	GoogleConfig *oauth2.Config
}

func NewAuthService() *AuthService {
	env := config.Env

	googleConfig := &oauth2.Config{
		ClientID:     env.GoogleID,
		ClientSecret: env.GoogleSecret,
		RedirectURL:  env.GoogleCallback,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	return &AuthService{
		GoogleConfig: googleConfig,
	}
}

type GoogleUser struct {
	Id            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"piture"`
	Locale        string `json:"locale"`
}
