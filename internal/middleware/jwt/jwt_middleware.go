package jwt_middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	error_response "github.com/reinaldo-silva/savina-stock/package/response/error"
)

type JwtMiddleware struct {
	secretKey []byte
}

type contextKey string

const (
	userIDKey   contextKey = "userID"
	userRoleKey contextKey = "userRole"
)

func NewJwtMiddleware(secretKey []byte) *JwtMiddleware {
	return &JwtMiddleware{secretKey: secretKey}
}

func (m *JwtMiddleware) ValidateToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			appError := error_response.NewAppError("Missing Authorization header", http.StatusUnauthorized)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(appError.StatusCode)
			json.NewEncoder(w).Encode(appError)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			appError := error_response.NewAppError("Invalid token format", http.StatusUnauthorized)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(appError.StatusCode)
			json.NewEncoder(w).Encode(appError)
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return m.secretKey, nil
		})

		if err != nil || !token.Valid {
			appError := error_response.NewAppError("Invalid or expired token", http.StatusUnauthorized)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(appError.StatusCode)
			json.NewEncoder(w).Encode(appError)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userID, _ := claims["user_id"].(float64)
			role, _ := claims["role"].(string)

			ctx := context.WithValue(r.Context(), userIDKey, uint(userID))
			ctx = context.WithValue(ctx, userRoleKey, role)
			r = r.WithContext(ctx)
		} else {
			appError := error_response.NewAppError("Could not parse token claims", http.StatusUnauthorized)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(appError.StatusCode)
			json.NewEncoder(w).Encode(appError)
			return
		}

		next.ServeHTTP(w, r)
	})
}
