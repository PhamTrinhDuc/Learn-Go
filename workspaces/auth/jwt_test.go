package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"testing"

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
