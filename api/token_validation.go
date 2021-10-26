package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
)

// validateToken will balidate bearer authorization header for protected routes
func (a *APIServer) validateToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")

		// extract bearer authorization token from request header
		bearerAuthHeader := r.Header.Get("authorization")

		token, httpStatus, err := a.authorizationHeaderValidation(bearerAuthHeader, " ", "Bearer")

		if err != nil {
			w.WriteHeader(httpStatus)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"result":  false,
				"message": fmt.Sprintf("%v", err),
			})
			return
		}

		// validate token if there is not error
		if token.Valid {
			context.Set(r, "token", token.Claims)

			// MJ Notes:
			//
			// DO some stuff to save token accesses on database on the future
			//
			//

			next.ServeHTTP(w, r)
		}
	})
}

// validateTokenWebSocket validate
func (a *APIServer) validateTokenWebSocket(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// extract bearer authorization token from request header
		jwtToken := r.Header.Get("Sec-WebSocket-Protocol")
		token, httpStatus, err := a.authorizationHeaderValidation(jwtToken, ", ", "token")

		if err != nil {
			w.WriteHeader(httpStatus)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"result":  false,
				"message": fmt.Sprintf("%v", err),
			})
			return
		}

		// validate token if there is not error
		if token.Valid {
			context.Set(r, "token_jwt", token.Claims)

			// MJ Notes:
			//
			// DO some stuff to save token accesses on database for future usage
			//
			//

			next.ServeHTTP(w, r)
		}
	})
}

// authorizationHeaderValidation checks the auth header and return token, http status code and error
func (a *APIServer) authorizationHeaderValidation(bearerAuthHeader, splitChar, firstWord string) (*jwt.Token, int, error) {
	// check if bearer auth header is set
	if bearerAuthHeader == "" {
		// header not set
		return nil, http.StatusBadRequest, errors.New("Token is not provided in authorization header")
	}

	// split bearer and jwt token
	authToken := strings.Split(bearerAuthHeader, splitChar)
	if len(authToken) != 2 {
		// header set but not in a correct way
		return nil, http.StatusBadRequest, errors.New("Authorization header provided but not in the correct format. correct format Authorization: Bearer {token-string}")
	}

	// check if string first array is firstWord and second is a VALID TOKEN
	if authToken[0] != firstWord {
		// header set but not in a correct way
		return nil, http.StatusBadRequest, errors.New("Authorization header provided but not in the correct format. correct format Authorization: Bearer {token-string}")
	}

	// Check token string
	token, err := jwt.Parse(authToken[1], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Invalid Token")
		}
		return []byte(a.apiConf.JWTSecretToken), nil
	})
	// Check for parsing error
	if err != nil {
		return nil, http.StatusUnauthorized, errors.New("Provided token is not valid! Please login and provide a valid token")
	}

	return token, http.StatusOK, nil
}
