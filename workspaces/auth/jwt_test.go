package auth

import (
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
				tokenString, _ := token.SignedString(token)
				return tokenString
			},
			wantErr:     true,
			errContains: "expired",
		},
		// 4. test case 4
		{
			name: "wrong issuer",
			tokenFunc: func() string {
				claims := Claims{
					TenantID: "tenant_id-123",
					UserID:   "user_id-123",
					Scopes:   []string{"write", "read"},
					RegisteredClaims: jwt.RegisteredClaims{
						Issuer:   "wrong issuer", // khác với issuer của validator phía trên,
						Audience: jwt.ClaimStrings{"mcp-demo"},
					},
				}
				token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
				tokenString, _ := token.SignedString(token)
				return tokenString
			},
			wantErr:     true,
			errContains: "wrong issuer",
		},

		// 5. test case 5
		{
			name: "wrong audience",
			tokenFunc: func() string {

			},
		},

		// 6. test case 6
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
