package middleware

import (
	"auth"
	"encoding/json"
	"net/http"
	protocol "protocal"
)

type AuthMiddle struct {
	validator           *auth.JWTValidator
	allowUnautheticated map[string]bool // cho phép các endpoint không cần auth
}

// NewAuthMiddleware creates a new auth middleware
func NewAuthMiddleware(validator *auth.JWTValidator) *AuthMiddle {
	return &AuthMiddle{
		validator: validator,
		allowUnautheticated: map[string]bool{
			protocol.MethodInitialize: true,
		},
	}
}

func (am *AuthMiddle) SendError(w http.ResponseWriter, id interface{}, code int, message string) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)

	response := protocol.NewErrorResponse(id, code, message, nil)
	json.NewEncoder(w).Encode(response)
}

// Handler wraps an HTTP handler with authentication
func (am *AuthMiddle) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Get header
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			am.SendError(w, nil, protocol.AuthenticationRequired, "Authorization header required")
			return
		}

		// 2. Validation token
		claims, err := am.validator.ValidateToken(authHeader)
		if err != nil {
			am.SendError(w, nil, protocol.AuthenticationRequired, "Invalid token: "+err.Error())
			return
		}

		// 3. Add auth to conttext
		ctx := auth.WithAuth(r.Context(), claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (am *AuthMiddle) OptionalHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			claims, err := am.validator.ValidateToken(authHeader)
			if err != nil {
				am.SendError(w, nil, protocol.AuthenticationRequired, "Invalid token: "+err.Error())
				return
			}
			ctx := auth.WithAuth(r.Context(), claims)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		next.ServeHTTP(w, r)
	})
}
