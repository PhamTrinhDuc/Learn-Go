package main

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ========================= TEST DB CONNECTION ==================================
func TestNewDB(t *testing.T) {
	tests := []struct {
		name    string
		config  DBConfig
		wantErr bool
	}{
		// TC0: Missing 1 số trường quan trọng
		{
			name:    "missing host",
			config:  DBConfig{Port: 5432, DBName: "test"},
			wantErr: true,
		},
		{
			name:    "invalid port",
			config:  DBConfig{Host: "localhost", Port: 0, DBName: "test"},
			wantErr: true,
		},
		// TC1: pass thiếu 1 số trường
		{
			name: "Parse config err",
			config: DBConfig{
				Host:     "localhost",
				Port:     5433,
				User:     "postgres",
				Password: "postgres",
				DBName:   "vchatbot",
			},
			wantErr: false,
		},
		// TC2: Pool connection & Ping (on db => false/ off db => true)
		{
			name: "Pool connection err",
			config: DBConfig{
				Host:     "localhost",
				Port:     5433,
				User:     "postgres",
				Password: "postgres",
				DBName:   "vchatbot",
				SSLMode:  "disable",
				MaxConns: 10,
				MinConns: 2,
			},
			wantErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, err := NewDB(context.Background(), tc.config)
			t.Log(err)
			if tc.wantErr {
				assert.Error(t, err)
				assert.Nil(t, db)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, db)
			}
		})
	}
}

// ========================= TEST CÁC FUNCTION KHI DB ĐÃ READY ==================================

var TenantID string = "1111-1111111-111111111"

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func GetDBConfig() DBConfig {
	return DBConfig{
		Host:     getEnvOrDefault("DB_HOST", "localhost"),
		Port:     5433,
		User:     getEnvOrDefault("DB_USER", "postgres"),
		Password: getEnvOrDefault("DB_PASSWORD", "postgres"),
		DBName:   getEnvOrDefault("PASSWORD", "vchatbot"),
		SSLMode:  getEnvOrDefault("DB_SSLMODE", "disable"),
		MaxConns: 10,
		MinConns: 2,
	}
}

func SetupDB(t *testing.T) *DB {
	cfg := GetDBConfig()
	db, err := NewDB(context.Background(), cfg)
	require.NoError(t, err, "Failed to connect to test database")
	return db
}

// integration test
func TestTxSetTenantID(t *testing.T) {
	db := SetupDB(t)
	_, err := db.BeginTx(context.Background(), TenantID)
	require.NoError(t, err, "Failed to set tenantID with transaction")
}
