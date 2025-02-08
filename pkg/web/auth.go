package web

import (
	"context"
	"errors"
	"github.com/womat/go-api-template/pkg/crypt"
	"github.com/womat/go-api-template/pkg/jwt_util"
	"net/http"
	"strings"
)

var (
	ErrUnauthorized = errors.New("not authorized")
)

type Config struct {
	ApiKey    crypt.EncryptedString
	JwtSecret crypt.EncryptedString
	JwtID     string
	AppName   string
}

// WithAuth is a middleware that checks if the request is authorized.
//   - If the request is not authorized, it returns a 401 Unauthorized response.
//   - If the request is authorized, it calls the next handler.
func WithAuth(h http.Handler, config Config) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {

			var user string
			authenticated := false

			if len(config.ApiKey.Value()) > 0 && checkApiKey(r, config.ApiKey) {
				authenticated = true
				user = "apikey"
			}

			if !authenticated && len(config.JwtSecret.Value()) > 0 && len(config.JwtID) > 0 {
				if claims, ok := checkJwtToken(r, config); ok {
					authenticated = true
					user = claims.User
				}
			}

			if !authenticated {
				Encode(w, http.StatusUnauthorized, NewApiError(ErrUnauthorized))
				return
			}

			// Store authenticated user in context
			ctx := context.WithValue(r.Context(), "user", user)

			// Pass modified request with updated context
			h.ServeHTTP(w, r.WithContext(ctx))
		},
	)
}

// checkApiKey checks if the request contains a valid API key.
//   - If the request does not contain a valid API key, it returns false.
//   - If the request contains a valid API key, it returns true.
//   - The API key is expected to be in the X-Api-Key header.
//   - The API key is compared to the given apiKey.
func checkApiKey(r *http.Request, apiKey crypt.EncryptedString) bool {
	if key := r.Header.Get("X-Api-Key"); key != "" && key == apiKey.Value() {
		return true
	}
	return false
}

// checkJwtToken checks if the request contains a valid JWT token.
//   - If the request does not contain a valid JWT token, it returns false.
//   - If the request contains a valid JWT token, it returns true.
//   - The JWT token is expected to be in the Authorization header.
//   - The JWT token is validated using the given JWT secret and JWT ID.
//   - The claims of the JWT token are returned if the token is valid.
func checkJwtToken(r *http.Request, config Config) (*jwt_util.Claims, bool) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, false
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return nil, false
	}

	claims, err := jwt_util.ValidateToken(parts[1], config.AppName, "auth", config.JwtID, config.JwtSecret.Value())
	if err != nil {
		return nil, false
	}
	return claims, true
}
