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

type Claims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
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

		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
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

		if claims, ok := token.Claims.(*Claims); ok && token.Valid {

			userID := claims.UserID
			role := claims.Role

			ctx := context.WithValue(r.Context(), userIDKey, userID)
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

func (m *JwtMiddleware) RequireRoles(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			role, ok := r.Context().Value(userRoleKey).(string)

			if !ok {
				appError := error_response.NewAppError("Role not found in context", http.StatusForbidden)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(appError.StatusCode)
				json.NewEncoder(w).Encode(appError)
				return
			}

			for _, allowedRole := range allowedRoles {
				if role == allowedRole {
					next.ServeHTTP(w, r)
					return
				}
			}

			appError := error_response.NewAppError("Access denied: insufficient permissions", http.StatusForbidden)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(appError.StatusCode)
			json.NewEncoder(w).Encode(appError)
		})
	}
}
