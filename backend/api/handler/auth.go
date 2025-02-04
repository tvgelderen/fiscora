package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/tvgelderen/fiscora/internal/auth"
	"github.com/tvgelderen/fiscora/internal/config"
	"github.com/tvgelderen/fiscora/internal/repository"
)

func (h *Handler) HandleOAuthLogin(c echo.Context) error {
	url := h.AuthService.GoogleConfig.AuthCodeURL(config.Env.SessionSecret)
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *Handler) HandleOAuthCallback(c echo.Context) error {
	logger := getLogger(c)
	query := c.Request().URL.Query()
	state := query.Get("state")
	if state != config.Env.SessionSecret {
		return c.String(http.StatusForbidden, "Invalid state")
	}

	error := query.Get("error")
	if error != "" {
		return c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s/login", config.Env.FrontendUrl))
	}

	code := query.Get("code")
	token, err := h.AuthService.GoogleConfig.Exchange(context.Background(), code)
	if err != nil {
		logger.Error(fmt.Sprintf("Error getting user token: %v", err.Error()))
		return c.String(http.StatusInternalServerError, "Something went wrong")
	}

	client := h.AuthService.GoogleConfig.Client(context.Background(), token)
	response, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		logger.Error(fmt.Sprintf("Error getting user info: %v", err.Error()))
		return c.String(http.StatusInternalServerError, "Something went wrong")
	}

	var googleUser auth.GoogleUser
	if err = json.NewDecoder(response.Body).Decode(&googleUser); err != nil {
		logger.Error(fmt.Sprintf("Error decoding user info: %v", err.Error()))
		return c.String(http.StatusInternalServerError, "Something went wrong")
	}

	user, err := h.UserRepository.GetByProviderId(c.Request().Context(), "google", googleUser.Id)
	if err != nil {
		if !repository.NoRowsFound(err) {
			logger.Error(fmt.Sprintf("Error getting user from db: %v", err.Error()))
			return c.String(http.StatusInternalServerError, "Something went wrong")
		}

		user, err = h.UserRepository.Add(c.Request().Context(), repository.CreateUserParams{
			ID:         uuid.New(),
			Provider:   "google",
			ProviderID: googleUser.Id,
			Username:   googleUser.Name,
			Email:      googleUser.Email,
		})
		if err != nil {
			logger.Error(fmt.Sprintf("Error creating user: %v", err.Error()))
			return c.String(http.StatusInternalServerError, "Something went wrong")
		}
	}

	authToken, err := auth.CreateToken(user.ID, user.Username, user.Email)
	if err != nil {
		logger.Error(fmt.Sprintf("Error creating auth token: %v", err.Error()))
		return c.String(http.StatusInternalServerError, "Something went wrong")
	}

	auth.SetToken(c.Response().Writer, authToken)

	return c.Redirect(http.StatusTemporaryRedirect, config.Env.FrontendUrl)
}

func (h *Handler) HandleDemoLogin(c echo.Context) error {
	logger := getLogger(c)
	demo, err := h.UserRepository.GetByEmail(c.Request().Context(), "demo")
	if err != nil {
		logger.Error(fmt.Sprintf("Error getting demo user from db: %v", err.Error()))
		return c.String(http.StatusInternalServerError, "Something went wrong")
	}

	authToken, err := auth.CreateToken(demo.ID, demo.Username, demo.Email)
	if err != nil {
		logger.Error(fmt.Sprintf("Error creating auth token for demo user: %v", err.Error()))
		return c.String(http.StatusInternalServerError, "Something went wrong")
	}

	auth.SetToken(c.Response().Writer, authToken)

	return c.Redirect(http.StatusTemporaryRedirect, config.Env.FrontendUrl)
}

func (h *Handler) HandleLogout(c echo.Context) error {
	auth.DeleteToken(c.Response().Writer)

	return c.Redirect(http.StatusTemporaryRedirect, config.Env.FrontendUrl)
}
