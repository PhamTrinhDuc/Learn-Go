package auth

import (
	"crypto/rsa"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateDemoTokenWithExpiry(
	tenantID string,
	userID string,
	scopes []string,
	privateKey *rsa.PrivateKey,
	expiry time.Duration,
) (string, error) {
	now := time.Now()
	claims := Claims{
		TenantID: tenantID,
		UserID:   userID,
		Scopes:   scopes,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "mcp-server-demo",
			Audience:  jwt.ClaimStrings{"mcp-server"},
			ExpiresAt: jwt.NewNumericDate(now.Add(expiry)), // thời điểm hết hạn
			IssuedAt:  jwt.NewNumericDate(now),             // thời điểm phát hành
			NotBefore: jwt.NewNumericDate(now),             //
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}
	return tokenString, nil
}

func GenerateDemoToken(tenantID, userID string, scopes []string, privateKey *rsa.PrivateKey) (string, error) {
	return GenerateDemoTokenWithExpiry(tenantID, userID, scopes, privateKey, 24*time.Hour)
}
