package auth

import (
	"context"
	"crypto/rsa"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type ContextKey string

const (
	// context key for tenant ID
	ContextKeyTenantID ContextKey = "tenant_id"
	// context key for user ID
	ContextKeyUserID ContextKey = "user_id"
	// context key for authorization scopes
	ContextKeyScopes ContextKey = "scopes"
)

type JWTValidator struct {
	publicKey *rsa.PublicKey
	issuer    string
	audience  string
}

// Cấu hình cho JWT
type Config struct {
	PublicKeyPEM string // RSA public key in PEM format
	Issuer       string // Expected token issuer
	Audience     string // Expected token audience
}

// Cấu trúc để JWTValidator parse token ra
type Claims struct {
	TenantID string   `json: "tenant_id"` // tag json: ánh xạ field trong struct sang key trong json
	UserID   string   `json: "user_id"`
	Email    string   `json: "email,omitempty"`
	Scopes   []string `json: "scopes,omitempty"`
	jwt.RegisteredClaims
}

// Khởi tạo new JWT Validator (constructor)
func NewJWTValidator(cfg Config) (*JWTValidator, error) {
	// 1. Parse RSA Public key from PEM
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(cfg.PublicKeyPEM))
	if err != nil {
		return nil, fmt.Errorf("failed to parse Publickey in format PEM")
	}
	// 2. Return address JWTValidator for pointer
	return &JWTValidator{
		publicKey: publicKey,
		issuer:    cfg.Issuer,
		audience:  cfg.Audience,
	}, nil
}

// ValidateToken validates JWT token and return claims
func (v *JWTValidator) ValidateToken(tokenString string) (*Claims, error) {
	// 1.Remove Bearer prefix in tokenString
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	// 2. Parse and validate token
	token, err := jwt.ParseWithClaims(
		tokenString, // token để parse
		&Claims{},   // khuôn để dổ dữ liệu vào (pointer để có thể ghi dữ liệu vào)
		func(token *jwt.Token) (interface{}, error) { // callback để cấp key verify
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok { // token.Method có kiểu jwt.SigningMethod - 1 interface. Go chưa biết là RSA hay HMAC. Dùng .(type) để hỏi: token.Method có phải là *jwt.SigningMethodRSA không
				return nil, fmt.Errorf("unexprected signing method: %v", token.Header["alg"])
			}
			return v.publicKey, nil
		})
	if err != nil { // không phải kiểu *jwt.SigningMethodRSA
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}
	// có phải Claims có các trường như đã định nghĩa không
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("Invalid token claims")
	}
	// 3. Validate Issuer and Audience
	if claims.Issuer != v.issuer {
		return nil, fmt.Errorf("Invalid issuer: expected %s, got %s", v.issuer, claims.Issuer)
	}
	validAudience := false
	for _, aud := range claims.Audience {
		if aud == v.audience {
			validAudience = true
			break
		}
	}
	if !validAudience {
		return nil, fmt.Errorf("Invalid audience")
	}

	// 4. Validate expriration
	if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
		return nil, fmt.Errorf("token expired")
	}

	// 5. Validate tenant ID is present
	if claims.TenantID == "" {
		return nil, fmt.Errorf("tenant_id in claim is required")
	}
	return claims, nil

}

func ExtractTenantID(ctx context.Context) (string, error) {
	tenantID, oke := ctx.Value(ContextKeyTenantID).(string)
	if !oke || tenantID == "" {
		return "", fmt.Errorf("tanent_id not found in context")
	}
	return tenantID, nil
}

func WithAuth(ctx context.Context, claims *Claims) context.Context {
	ctx = context.WithValue(ctx, ContextKeyUserID, claims.UserID)
	ctx = context.WithValue(ctx, ContextKeyTenantID, claims.TenantID)
	ctx = context.WithValue(ctx, ContextKeyScopes, claims.Scopes)
	return ctx
}

func ExtractUserID(ctx context.Context) (string, error) {
	userID, oke := ctx.Value(ContextKeyUserID).(string)
	if !oke || userID == "" {
		return "", fmt.Errorf("user_id not found in context")
	}
	return userID, nil
}

func ExtractScopes(ctx context.Context) ([]string, error) {
	scopes, oke := ctx.Value(ContextKeyScopes).([]string)
	if !oke || len(scopes) == 0 {
		return []string{}, fmt.Errorf("scopes not found in context")
	}
	return scopes, nil
}

func hasScope(ctx context.Context, requiredScope string) bool {
	scopes, err := ExtractScopes(ctx)
	if err != nil {
		return false
	}
	for _, scope := range scopes {
		if scope == requiredScope {
			return true
		}
	}
	return false
}
