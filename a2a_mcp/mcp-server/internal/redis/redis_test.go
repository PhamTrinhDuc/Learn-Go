package redis

import (
	"testing"

	utils "learn-go/a2a_mcp/mcp-server/internal/utils"

	"github.com/stretchr/testify/assert"
)

// ========================= TEST DB CONNECTION ==================================
func TestNewRedis(t *testing.T) {
	env := GetRedisConfig()
	tests := []struct {
		name    string
		config  RedisConfig
		wantErr bool
	}{
		// TC0: Missing 1 số trường quan trọng
		{
			name:    "missing host",
			config:  RedisConfig{Port: 6379, DB: 5, Password: env.Password, Username: env.Username},
			wantErr: true,
		},
		{
			name:    "invalid port",
			config:  RedisConfig{Host: env.Host, Port: 0, Password: env.Password, Username: env.Username},
			wantErr: true,
		},
		{
			name:    "invalid username",
			config:  RedisConfig{Host: env.Host, Port: env.Port, Username: "", Password: env.Password},
			wantErr: true,
		},
		{
			name:    "invalid password",
			config:  RedisConfig{Host: env.Host, Port: env.Port, Username: env.Username, Password: ""},
			wantErr: true,
		},
		{
			name:    "connect success",
			config:  RedisConfig{Host: env.Host, Port: env.Port, Username: env.Username, Password: env.Password},
			wantErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.config.Validate()
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func GetRedisConfig() RedisConfig {
	return RedisConfig{
		Host:     utils.GetEnvOrDefault("REDIS_HOST", "localhost"),
		Port:     6379,
		Username: utils.GetEnvOrDefault("REDIS_USERNAME", "jiyuu"),
		Password: utils.GetEnvOrDefault("REDIS_PASSWORD", "a2amcpgo"),
	}
}
