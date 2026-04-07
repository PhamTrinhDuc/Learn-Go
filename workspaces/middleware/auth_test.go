package middleware

import (
	"auth"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"middleware"
	protocol "protocal"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setTestAuth(t *testing.T) (*auth.JWTValidator, *rsa.PrivateKey, string) {
	// 1. Generate private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	// 2. Generate public key
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	require.NoError(t, err)
	// 3. Generate publicKey format PEM
	publickeyPEM := string(pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC_KEY",
		Bytes: publicKeyBytes,
	}))

	// 4. Create authValidator
	validator, err := auth.NewJWTValidator(
		auth.Config{
			PublicKeyPEM: publickeyPEM,
			Issuer:       "mcp-server-demo",
			Audience:     "mcp-server",
		},
	)
	require.NoError(t, err)
	return validator, privateKey, publickeyPEM
}

func TestNewMiddleware(t *testing.T) *middleware.AuthMiddle {
	validator, _, _ := setTestAuth(t)
	middleware := NewAuthMiddleware(validator)

	assert.NotNil(t, middleware)
	assert.NotNil(t, middleware.validator)
	assert.NotNil(t, middleware.allowUnautheticated)
	assert.True(t, middleware.allowUnautheticated[protocol.MethodInitialize])
	return middleware
}

func TestMiddlewareHandleValidToken(t *testing.T) {
	middleware := TestNewMiddleware(t)
}
