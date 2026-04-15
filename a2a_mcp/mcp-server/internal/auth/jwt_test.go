package auth

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func generateKeyPair(t *testing.T) (*rsa.PrivateKey, string) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	require.NoError(t, err)

	publicKeyPEM := pem.EncodeToMemory((&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	}))

	return privateKey, string(publicKeyPEM)
}
func TestNewJWTValidator(t *testing.T) {
	_, publicKeyPEM := generateKeyPair(t)
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: Config{
				PublicKeyPEM: publicKeyPEM,
				Issuer:       "issuer-test",
				Audience:     "audience-test",
			},
			wantErr: false,
		},
		{
			name: "invalid public key PEM",
			config: Config{
				PublicKeyPEM: "invalid public key PEM",
				Issuer:       "issuer-test",
				Audience:     "audience-test",
			},
			wantErr: true,
		},
		// {
		// 	name: "missing field",
		// 	config: Config{
		// 		PublicKeyPEM: publicKeyPEM,
		// 		Issuer:       "issuer-test",
		// 	},
		// 	wantErr: true,
		// },
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			validator, err := NewJWTValidator(test.config)
			if test.wantErr {
				assert.Error(t, err)
				assert.Nil(t, validator)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, validator)
			}
		})
	}
}

func TestGenerateToken(t *testing.T) {
	privateKey, publicKeyPEM := generateKeyPair(t)

	validator, err := NewJWTValidator(Config{
		PublicKeyPEM: publicKeyPEM,
		Issuer:       "mcp-server-demo",
		Audience:     "mcp-server",
	})

	require.NoError(t, err)

	tests := []struct {
		name        string
		tokenFunc   func() string
		wantErr     bool
		errContains string
		validate    func(t *testing.T, claims *Claims)
	}{
		{
			// 1. test case 1
			name: "valid token",
			tokenFunc: func() string {
				token, _ := GenerateDemoToken(
					"tenant_id-123",
					"user_id-123",
					[]string{"read", "write"},
					privateKey,
				)
				return token
			},
			wantErr: false,
			validate: func(t *testing.T, claims *Claims) {
				require.Equal(t, "tenant_id-123", claims.TenantID)
				require.Equal(t, "user_id-123", claims.UserID)
				require.Contains(t, claims.Scopes, "read")
				require.Contains(t, claims.Scopes, "write")
			},
		},
		{
			// 2. test case 2
			name: "token with Bearer prefix",
			tokenFunc: func() string {
				token, _ := GenerateDemoToken(
					"tenant_id-123",
					"user_id-123",
					[]string{"read", "write"},
					privateKey,
				)
				return "Bearer" + token
			},
			wantErr:  false,
			validate: nil,
		},
		{
			// 3. test case 3
			name: "expired token",
			tokenFunc: func() string {
				now := time.Now()
				claims := Claims{
					TenantID: "tenant_id-123",
					UserID:   "user_id-123",
					Scopes:   []string{"read", "write"},
					RegisteredClaims: jwt.RegisteredClaims{
						Issuer:    "mcp-server-demo",
						Audience:  jwt.ClaimStrings{"mcp-server"},
						ExpiresAt: jwt.NewNumericDate(now.Add(-1 * time.Hour)), // hết hạn 1 tiếng trước
						IssuedAt:  jwt.NewNumericDate(now.Add(-2 * time.Hour)), // tạo 2 tiếng trước
					},
				}
				token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
				tokenString, _ := token.SignedString(privateKey)
				return tokenString
			},
			wantErr:     true,
			errContains: "expired",
		},
		// 4. test case 4
		{
			name: "wrong issuer",
			tokenFunc: func() string {
				now := time.Now()
				claims := Claims{
					TenantID: "tenant_id-123",
					UserID:   "user_id-123",
					Scopes:   []string{"write", "read"},
					RegisteredClaims: jwt.RegisteredClaims{
						Issuer:    "wrong issuer", // khác với issuer của validator phía trên,
						Audience:  jwt.ClaimStrings{"mcp-demo"},
						ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour)),
						IssuedAt:  jwt.NewNumericDate(now),
					},
				}
				token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
				tokenString, _ := token.SignedString(privateKey)
				return tokenString
			},
			wantErr:     true,
			errContains: "Invalid issuer",
		},

		// 5. test case 5
		{
			name: "wrong audience",
			tokenFunc: func() string {
				now := time.Now()
				claims := Claims{
					TenantID: "tenant_id-123",
					UserID:   "user_id-123",
					Scopes:   []string{"write", "read"},
					RegisteredClaims: jwt.RegisteredClaims{
						Issuer:    "mcp-server-demo",
						Audience:  jwt.ClaimStrings{"wrong-audience"}, // sai so với validator
						ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour)),
						IssuedAt:  jwt.NewNumericDate(now),
					},
				}
				token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
				tokenString, _ := token.SignedString(privateKey)
				return tokenString
			},
			wantErr:     true,
			errContains: "Invalid audience",
		},

		// 6. test case 6
		{
			name: "wrong signing method",
			tokenFunc: func() string {
				now := time.Now()
				claims := Claims{
					TenantID: "tenant_id-123",
					UserID:   "user_id-123",
					Scopes:   []string{"write", "read"},
					RegisteredClaims: jwt.RegisteredClaims{
						Issuer:    "mcp-server-demo",
						Audience:  jwt.ClaimStrings{"mcp-server"},
						ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour)),
						IssuedAt:  jwt.NewNumericDate(now),
					},
				}
				// Dùng HS256 với secret là byte array thay vì RSA private key để pass qua SignedString
				// nhưng sẽ fail ở validator vì validator check RSA
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				tokenString, _ := token.SignedString([]byte("secret"))
				return tokenString
			},
			wantErr:     true,
			errContains: "unexpected signing method",
		},
		{
			name: "missing tenant id",
			tokenFunc: func() string {
				now := time.Now()
				claims := Claims{
					UserID: "user_id-123",
					Scopes: []string{"write", "read"},
					RegisteredClaims: jwt.RegisteredClaims{
						Issuer:    "mcp-server-demo",
						Audience:  jwt.ClaimStrings{"mcp-server"},
						ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour)),
						IssuedAt:  jwt.NewNumericDate(now),
					},
				}
				token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
				tokenString, _ := token.SignedString(privateKey)
				return tokenString
			},
			wantErr:     true,
			errContains: "tenant_id in claim is required",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tokenString := tc.tokenFunc()
			claims, err := validator.ValidateToken(tokenString)
			if tc.wantErr {
				require.Error(t, err)
				if tc.errContains != "" {
					require.Contains(t, err.Error(), tc.errContains)
				}
			} else {
				require.NoError(t, err)
				require.NotNil(t, claims)
				if tc.validate != nil {
					tc.validate(t, claims)
				}
			}
		})
	}
}

func TestExtractTenantID(t *testing.T) {
	tests := []struct {
		name     string
		ctx      context.Context
		expected string
		wantErr  bool
	}{
		{
			name:     "valid tenant ID",
			ctx:      context.WithValue(context.Background(), ContextKeyTenantID, "tenant_id-123"),
			expected: "tenant_id-123",
			wantErr:  false,
		},
		{
			name:    "missing tenant ID",
			ctx:     context.Background(),
			wantErr: true,
		},
		{
			name:     "invalid tenant ID",
			ctx:      context.WithValue(context.Background(), ContextKeyTenantID, 123), // using tenantID format integer to fail context.Value()
			expected: "tenant_id-123",
			wantErr:  true,
		},
		{
			name:     "empty tenant ID",
			ctx:      context.WithValue(context.Background(), ContextKeyTenantID, ""),
			expected: "tenant_id-123",
			wantErr:  true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tenantID, err := ExtractTenantID(tc.ctx)
			if tc.wantErr {
				require.Error(t, err)
				// require.Nil(t, tenantID)
			} else {
				require.NoError(t, err)
				require.NotNil(t, tenantID)
				require.Equal(t, tc.expected, tenantID)
			}
		})
	}
}

func TestExtractUserID(t *testing.T) {
	tests := []struct {
		name     string
		ctx      context.Context
		expected string
		wantErr  bool
	}{
		{
			name:     "valid user ID",
			ctx:      context.WithValue(context.Background(), ContextKeyUserID, "user_id-123"),
			expected: "user_id-123",
			wantErr:  false,
		},
		{
			name:    "missing user ID",
			ctx:     context.Background(),
			wantErr: true,
		},
		{
			name:     "invalid user ID",
			ctx:      context.WithValue(context.Background(), ContextKeyUserID, 123), // using tenantID format integer to fail context.Value()
			expected: "user_id-123",
			wantErr:  true,
		},
		{
			name:     "empty user ID",
			ctx:      context.WithValue(context.Background(), ContextKeyUserID, ""),
			expected: "user_id-123",
			wantErr:  true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			userID, err := ExtractUserID(tc.ctx)
			if tc.wantErr {
				require.Error(t, err)
				// require.Nil(t, userID)
			} else {
				require.NoError(t, err)
				require.NotNil(t, userID)
				require.Equal(t, tc.expected, userID)
			}
		})
	}
}

func TestExtractScopes(t *testing.T) {
	tests := []struct {
		name     string
		ctx      context.Context
		expected []string
		wantErr  bool
	}{
		{
			name:     "valid Scopes",
			ctx:      context.WithValue(context.Background(), ContextKeyScopes, []string{"write", "read"}),
			expected: []string{"write", "read"},
			wantErr:  false,
		},
		{
			name:    "missing Scope",
			ctx:     context.Background(),
			wantErr: true,
		},
		{
			name:     "invalid Scopes",
			ctx:      context.WithValue(context.Background(), ContextKeyScopes, 123),
			expected: []string{"write", "read"},
			wantErr:  true,
		},
		{
			name:     "empty Scopes",
			ctx:      context.WithValue(context.Background(), ContextKeyScopes, []string{}),
			expected: []string{"write", "read"},
			wantErr:  true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			Scopes, err := ExtractScopes(tc.ctx)
			if tc.wantErr {
				require.Error(t, err)
				// require.Nil(t, Scopes)
			} else {
				require.NoError(t, err)
				require.NotNil(t, Scopes)
				require.Equal(t, tc.expected, Scopes)
			}
		})
	}
}

func TestHasScope(t *testing.T) {
	tests := []struct {
		name          string
		requiredScope string
		ctx           context.Context
		wantErr       bool
	}{
		{
			name:          "valid scopes",
			requiredScope: "read",
			ctx:           context.WithValue(context.Background(), ContextKeyScopes, []string{"write", "read"}),
			wantErr:       true,
		},
		{
			name:          "invalid scopes",
			requiredScope: "delete",
			ctx:           context.WithValue(context.Background(), ContextKeyScopes, []string{"write", "read"}),
			wantErr:       false,
		},
		{
			name:          "empty scopes",
			requiredScope: "delete",
			ctx:           context.Background(),
			wantErr:       false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			flag := hasScope(tc.ctx, tc.requiredScope)
			require.Equal(t, flag, tc.wantErr)
		})
	}
}

func TestWithAuth(t *testing.T) {
	claims := Claims{
		TenantID: "tenant_id-123",
		UserID:   "user_id-123",
		Scopes:   []string{"read", "write"},
	}
	ctx := context.Background()

	ctx = WithAuth(ctx, &claims)
	// validate values equation
	tenantID, err := ExtractTenantID(ctx)
	require.NoError(t, err)
	require.Equal(t, claims.TenantID, tenantID)

	userID, err := ExtractUserID(ctx)
	require.NoError(t, err)
	require.Equal(t, claims.UserID, userID)

	scopes, err := ExtractScopes(ctx)
	require.NoError(t, err)
	require.Equal(t, claims.Scopes, scopes)
}

// go test -v .
