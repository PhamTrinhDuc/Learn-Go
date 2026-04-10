package middleware

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"learn-go/workspaces/auth"
	"learn-go/workspaces/protocol"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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

func TestNewMiddleware(t *testing.T) {
	validator, _, _ := setupTestAuth(t)
	middleware := NewAuthMiddleware(validator)

	assert.NotNil(t, middleware)
	assert.NotNil(t, middleware.validator)
	assert.NotNil(t, middleware.allowUnautheticated)
	assert.True(t, middleware.allowUnautheticated[protocol.MethodInitialize])
}

func TestMiddleware_Handle_ValidToken(t *testing.T) {
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
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder() // lưu trữ dữ liệu response vào đây thay cho client/trình duyệt của client

	// Execute
	handler := middleware.Handler(testHandler)
	handler.ServeHTTP(rr, req)

	require.Equal(t, handleCalled, true)
	require.Equal(t, http.StatusOK, rr.Code)
}

func TestMiddleware_Handle_MissingToken(t *testing.T) {
	validator, _, _ := setupTestAuth(t)
	middleware := NewAuthMiddleware(validator)
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("Handler should not be called missing token")
	})

	// Create request and store response
	req := httptest.NewRequest("POST", "/mcp", nil)
	rr := httptest.NewRecorder()

	// Execute
	handler := middleware.Handler(testHandler)
	handler.ServeHTTP(rr, req)

	// Verify err code
	require.Equal(t, http.StatusUnauthorized, rr.Code)

	var response protocol.Response
	err := json.NewDecoder(rr.Body).Decode(&response)
	require.NoError(t, err)
	assert.NotNil(t, response.Error)
	assert.NotNil(t, protocol.AuthenticationRequired, response.Error.Code)
	assert.Contains(t, response.Error.Message, "Authorization header required")
}

func TestMiddleware_Handle_InValidToken(t *testing.T) {
	validator, _, _ := setupTestAuth(t)
	middleware := NewAuthMiddleware(validator)

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("Handler should not be called invalid token")
	})

	// create request and store response
	req := httptest.NewRequest("POST", "/mcp", nil)
	req.Header.Set("Authorization", "Invalid token")
	rr := httptest.NewRecorder()

	// Execute
	hanlder := middleware.Handler(testHandler)
	hanlder.ServeHTTP(rr, req)

	// Verify err response
	require.Equal(t, http.StatusUnauthorized, rr.Code)

	var response protocol.Response

	err := json.NewDecoder(rr.Body).Decode(&response)
	require.NoError(t, err)
	assert.NotNil(t, response.Error)
	assert.Equal(t, protocol.AuthenticationRequired, response.Error.Code)
	assert.Contains(t, response.Error.Message, "Invalid token")
}

func TestMiddleware_Handle_ExpiredToken(t *testing.T) {
	validator, privateKey, _ := setupTestAuth(t)
	token, err := auth.GenerateDemoTokenWithExpiry(
		"tenant_id-123",
		"user_id-123",
		[]string{"admin"},
		privateKey,
		-time.Hour,
	)
	require.NoError(t, err)
	middleware := NewAuthMiddleware(validator)

	// Create test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("Handler should not be called expired token")
	})

	// create request and  object store response
	req := httptest.NewRequest("POST", "/mcp", nil)
	req.Header.Set("Authorization", "Bearer"+token)
	rr := httptest.NewRecorder() // lưu trữ dữ liệu response vào đây thay cho client/trình duyệt của client

	// Execute
	hanlder := middleware.Handler(testHandler)
	hanlder.ServeHTTP(rr, req)

	// Verify err response
	require.Equal(t, http.StatusUnauthorized, rr.Code)

	var response protocol.Response

	errRes := json.NewDecoder(rr.Body).Decode(&response)
	require.NoError(t, errRes)
	assert.NotNil(t, response.Error)
	assert.Equal(t, protocol.AuthenticationRequired, response.Error.Code)
	assert.Contains(t, response.Error.Message, "failed to parse token") // failed in JWTValidator.ValidateToken
}

func TestMiddleware_OptionHandle_ValidToken(t *testing.T) {
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
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder() // lưu trữ dữ liệu response vào đây thay cho client/trình duyệt của client

	// Execute
	handler := middleware.OptionalHandler(testHandler)
	handler.ServeHTTP(rr, req)

	require.Equal(t, handleCalled, true)
	require.Equal(t, http.StatusOK, rr.Code)
}

func TestMiddleware_OptionHandle_MissingToken(t *testing.T) {
	validator, _, _ := setupTestAuth(t)
	middleware := NewAuthMiddleware(validator)
	handleCalled := false
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handleCalled = true

		// tại sao user_id và tenant_id not found in context?
		// vì không tạo token (có chứa user_id và tenant_id) như các hàm Valid.
		// Do đó khi không truyền token có chứa 2 trường này => lỗi not found

		// tenantID, err := auth.ExtractTenantID(r.Context())
		// require.NoError(t, err)
		// require.Equal(t, "tenant_id-123", tenantID)

		// userID, err := auth.ExtractUserID(r.Context())
		// require.NoError(t, err)
		// require.Equal(t, "user_id-123", userID)

		w.WriteHeader(http.StatusOK)
	})

	// Create request and store response
	req := httptest.NewRequest("POST", "/mcp", nil)
	rr := httptest.NewRecorder()

	// Execute
	handler := middleware.OptionalHandler(testHandler)
	handler.ServeHTTP(rr, req)

	// Verify err code
	assert.Equal(t, handleCalled, true)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestMiddleware_OptionHandle_InValidToken(t *testing.T) {
	validator, _, _ := setupTestAuth(t)
	middleware := NewAuthMiddleware(validator)

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("Handler should not be called invalid token")
	})

	// create request and store response
	req := httptest.NewRequest("POST", "/mcp", nil)
	req.Header.Set("Authorization", "Invalid token")
	rr := httptest.NewRecorder()

	// Execute
	hanlder := middleware.OptionalHandler(testHandler)
	hanlder.ServeHTTP(rr, req)

	// Verify err response
	require.Equal(t, http.StatusUnauthorized, rr.Code)

	var response protocol.Response

	err := json.NewDecoder(rr.Body).Decode(&response)
	require.NoError(t, err)
	assert.NotNil(t, response.Error)
	assert.Equal(t, protocol.AuthenticationRequired, response.Error.Code)
	assert.Contains(t, response.Error.Message, "Invalid token")
}

func TestMiddleware_OptionHandle_ExpiredToken(t *testing.T) {
	validator, privateKey, _ := setupTestAuth(t)
	token, err := auth.GenerateDemoTokenWithExpiry(
		"tenant_id-123",
		"user_id-123",
		[]string{"admin"},
		privateKey,
		-time.Hour,
	)
	require.NoError(t, err)
	middleware := NewAuthMiddleware(validator)

	// Create test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("Handler should not be called expired token")
	})

	// create request and  object store response
	req := httptest.NewRequest("POST", "/mcp", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder() // lưu trữ dữ liệu response vào đây thay cho client/trình duyệt của client

	// Execute
	hanlder := middleware.OptionalHandler(testHandler)
	hanlder.ServeHTTP(rr, req)

	// Verify err response
	require.Equal(t, http.StatusUnauthorized, rr.Code)

	var response protocol.Response

	errRes := json.NewDecoder(rr.Body).Decode(&response)
	require.NoError(t, errRes)
	assert.NotNil(t, response.Error)
	assert.Equal(t, protocol.AuthenticationRequired, response.Error.Code)
	assert.Contains(t, response.Error.Message, "failed to parse token")
}
