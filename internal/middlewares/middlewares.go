package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Secret key để verify token
var jwtSecret = []byte("your-secret-key")

// Role-based access control (RBAC)
var rolePermissions = map[string][]string{
	"admin": {"GET", "POST", "PUT", "DELETE"},
	"user":  {"GET"},
}

// TrackRequestMiddleware adds an unique id to each request for it to be tracked
// afterwards in logs if needed.
func TrackRequestMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	ctx := r.Context()

	requestID := uuid.New().String()
	ctx = context.WithValue(ctx, "requestID", requestID)

	r = r.WithContext(ctx)

	next(w, r)
}

// Định nghĩa key riêng cho context
type contextKey string

const userKey contextKey = "user"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization Header", http.StatusUnauthorized)
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			http.Error(w, "Invalid Token Format", http.StatusUnauthorized)
			return
		}
		tokenString := tokenParts[1]

		claims := &jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Invalid Token", http.StatusUnauthorized)
			return
		}

		// Thêm user vào context với key riêng biệt
		ctx := context.WithValue(r.Context(), userKey, (*claims)["username"])
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

// Kiểm tra xem userRole có quyền thực hiện method không
func hasPermission(userRole, method string) bool {
	permissions, exists := rolePermissions[userRole]
	if !exists {
		return false
	}
	for _, perm := range permissions {
		if perm == method {
			return true
		}
	}
	return false
}
