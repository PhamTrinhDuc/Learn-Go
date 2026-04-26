package database

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
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
				User:     "mcp_user",
				Password: "mcp_password",
				DBName:   "mcp_db",
			},
			wantErr: false,
		},
		// TC2: Pool connection & Ping (on db => false/ off db => true)
		{
			name: "Pool connection err",
			config: DBConfig{
				Host:     "localhost",
				Port:     5433,
				User:     "mcp_user",
				Password: "mcp_password",
				DBName:   "mcp_db",
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


