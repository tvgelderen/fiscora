package auth

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/tvgelderen/fiscora/config"
)

const accessTokenKey = "AccessToken"

type CustomClaims struct {
	Id    uuid.UUID `json:"id"`
	Email string    `json:"email"`
	Name  string    `json:"name"`
	jwt.RegisteredClaims
}

func CreateToken(id uuid.UUID, name, email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, CustomClaims{
		Id:    id,
		Name:  name,
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	})

	return token.SignedString([]byte(config.Env.HMAC))
}

func SetToken(w http.ResponseWriter, token string) {
	cookie := http.Cookie{
		Name:     accessTokenKey,
		Value:    token,
		MaxAge:   36000,
		Path:     "/",
		HttpOnly: true,
		Secure:   config.Env.Production,
	}

	http.SetCookie(w, &cookie)
}

func DeleteToken(w http.ResponseWriter) {
	cookie := http.Cookie{
		Name:     accessTokenKey,
		Value:    "",
		MaxAge:   0,
		Path:     "/",
		HttpOnly: true,
		Secure:   config.Env.Production,
	}

	http.SetCookie(w, &cookie)
}

func GetId(r *http.Request) (uuid.UUID, error) {
	token, err := getToken(r)
	if err != nil {
		return uuid.UUID{}, err
	}

	parsedToken, err := parseToken(token)
	if err != nil {
		return uuid.UUID{}, err
	}

	claims := parsedToken.Claims.(*CustomClaims)
	if claims.ExpiresAt.Before(time.Now()) {
		return uuid.UUID{}, errors.New("token expired")
	}

	return claims.Id, nil
}

func getToken(r *http.Request) (string, error) {
	bearer := r.Header.Get("Authorization")
	if bearer != "" {
		splitBearer := strings.Split(bearer, " ")
		return splitBearer[len(splitBearer)-1], nil
	}

	cookie, err := r.Cookie(accessTokenKey)
	if err != nil {
		return "", err
	}

	return cookie.Value, nil
}

func parseToken(token string) (*jwt.Token, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Env.HMAC), nil
	})

	return parsedToken, err
}
