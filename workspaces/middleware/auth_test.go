package middleware

import (
	"auth"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"net/http"
	"net/http/httptest"
	protocol "protocal"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestAuth(t *testing.T) (*auth.JWTValidator, *rsa.PrivateKey, string) {
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

func TestNewMiddleware(t *testing.T) *AuthMiddle {
	validator, _, _ := setupTestAuth(t)
	middleware := NewAuthMiddleware(validator)

	assert.NotNil(t, middleware)
	assert.NotNil(t, middleware.validator)
	assert.NotNil(t, middleware.allowUnautheticated)
	assert.True(t, middleware.allowUnautheticated[protocol.MethodInitialize])
	return middleware
}

func TestMiddlewareHandleValidToken(t *testing.T) {
	validator, privateKey, _ := setupTestAuth(t)
	token, err := auth.GenerateDemoToken(
		"tenant_id-123",
		"user_id-123",
		[]string{"admin"},
		privateKey,
	)
	require.NoError(t, err)
	middleware := NewAuthMiddleware(validator)

	// Create test handler
	handleCalled := false
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handleCalled = true

		tenantID, err := auth.ExtractTenantID(r.Context())
		require.NoError(t, err)
		require.Equal(t, "tenant_id-123", tenantID)

		userID, err := auth.ExtractUserID(r.Context())
		require.NoError(t, err)
		require.Equal(t, "user_id-123", userID)

		w.WriteHeader(http.StatusOK)
	})

	// create request and  object store response
	req := httptest.NewRequest("POST", "/mcp", nil)
	req.Header.Set("Authorization", "Bearer"+token)
	rr := httptest.NewRecorder() // lưu trữ dữ liệu response vào đây thay cho client/trình duyệt của client

	// Execute
	handler := middleware.Handler(testHandler)
	handler.ServeHTTP(rr, req)

	require.Equal(t, handleCalled, true)
	require.Equal(t, http.StatusOK, rr.Code)
}

func TestMiddlewareHandleMissingToken(t *testing.T) {
	validator, _, _ := setupTestAuth(t)
	middleware := NewAuthMiddleware(validator)
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("Handlershould not be called invalid token")
	})

	// Create request and store response
	req := httptest.NewRequest("POST", "/mcp", nil)
	req.Header.Set("Authorization", "Invalid-token")
	rr := httptest.NewRecorder()

	// Execute
	handler := middleware.Handler(testHandler)
	handler.ServeHTTP(rr, req)

	// Verify
	require.Equal(t, http.StatusUnauthorized, rr.Code)

	var response protocol.Response
	err := json.NewDecoder(rr.Body).Decode(&response)
	require.NoError(t, err)
	assert.NotNil(t, err.Error())
	assert.NotNil(t, protocol.AuthenticationRequired, response.Error.Code)
	assert.Contains(t, response.Error.Message, "Authorizaton header required")

}
