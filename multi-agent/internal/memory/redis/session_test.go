package redis

import (
	"multi-agent/internal/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func GetRedisConfig() RedisConfig {
	return RedisConfig{
		Host:     utils.GetEnvString("REDIS_HOST", "localhost"),
		Port:     utils.GetEnvInt("REDIS_PORT", 6379),
		Username: utils.GetEnvString("REDIS_USERNAME", "jiyuu"),
		Password: utils.GetEnvString("REDIS_PASSWORD", "a2amcpgo"),
	}
}

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
			_, err := NewRedisService(&tc.config)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func SetupRedis() *RedisService {
	cfg := GetRedisConfig()
	client, err := NewRedisService(&cfg)
	return client
}

func TestCreateKeys(t *testing.T) {
	tests := []struct {
		name      string
		userID    string
		appName   string
		sessionID string
		expected  string
	}{
		{
			name:      "create app key",
			userID:    "",
			sessionID: "",
			appName:   "app01",
			expected:  "appstate:app01",
		},
		{
			name:      "create user key",
			userID:    "user001",
			sessionID: "",
			appName:   "app01",
			expected:  "users:app01:user001",
		},
		{
			name:      "create session key",
			userID:    "user001",
			appName:   "app01",
			sessionID: "123",
			expected:  "session:app01:user001:123",
		},
	}
	redisSvc := SetupRedis()
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var key string
			if tc.name == "create app key" {
				key = redisSvc.appStateKey(tc.appName)
			}
		})
	}
}
