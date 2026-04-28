package auth

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateToken(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret")
	defer os.Unsetenv("JWT_SECRET")

	t.Run("valid token", func(t *testing.T) {
		token, err := GenerateTokenWithPrivateKey("user123", "test@test.com", "admin")
		require.NoError(t, err)

		claims, err := ValidateToken(token)
		require.NoError(t, err)
		assert.Equal(t, "user123", claims.UserID)
		assert.Equal(t, "test@test.com", claims.Email)
		assert.Equal(t, "admin", claims.Role)
	})

	t.Run("expired token", func(t *testing.T) {
		token, err := GenerateTokenWithPrivateKey("user123", "test@test.com", "admin")
		require.NoError(t, err)

		_, err = ValidateToken(token)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "token is expired")
	})

	t.Run("invalid token string", func(t *testing.T) {
		_, err := ValidateToken("invalid-token")
		require.Error(t, err)
	})
}

func TestExtractTokenFromHeader(t *testing.T) {
	tests := []struct {
		name       string
		authHeader string
		want       string
		wantErr    bool
	}{
		{
			name:       "valid header",
			authHeader: "Bearer eyJhbG...",
			want:       "eyJhbG...",
			wantErr:    false,
		},
		{
			name:       "empty header",
			authHeader: "",
			wantErr:    true,
		},
		{
			name:       "missing bearer",
			authHeader: "eyJhbG...",
			wantErr:    true,
		},
		{
			name:       "invalid format",
			authHeader: "Bearer ",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExtractTokenFromHeader(tt.authHeader)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestGenerateDemoToken(t *testing.T) {
	token, err := GenerateTokenWithPrivateKey("user1", "user1@test.com", "user")
	require.NoError(t, err)
	require.NotEmpty(t, token)

	// verify token
	claims, err := ValidateToken(token)
	require.NoError(t, err)
	assert.Equal(t, "user1", claims.UserID)
}
