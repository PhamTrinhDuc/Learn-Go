package observability

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTelemtry(t *testing.T) {
	tests := []struct {
		name    string
		cfg     Config
		wantErr bool
	}{
		{
			name: "no enable services",
			cfg: Config{
				ServiceName:    "test-service",
				ServiceVersion: "version2",
				Enviroment:     "development",
				OTLPEndpoint:   "localhost:9999",
				SamplingRate:   0.3,
				EnableTracing:  false,
				EnableMetric:   false,
			},
			wantErr: false,
		},
		{
			name: "missing tracer",
			cfg: Config{
				ServiceName:    "test-service",
				ServiceVersion: "version2",
				Enviroment:     "development",
				OTLPEndpoint:   "localhost:9999",
				SamplingRate:   0.3,
				EnableTracing:  false,
				EnableMetric:   true,
			},
			wantErr: false,
		},
		{
			name: "missing meter",
			cfg: Config{
				ServiceName:    "test-service",
				ServiceVersion: "version2",
				Enviroment:     "development",
				OTLPEndpoint:   "localhost:9999",
				SamplingRate:   0.3,
				EnableTracing:  true,
				EnableMetric:   false,
			},
			wantErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			otel, err := NewTelemetry(context.Background(), tc.cfg)
			if tc.wantErr {
				assert.Error(t, err)
				assert.Nil(t, otel)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, otel)
			}

		})
	}
}
