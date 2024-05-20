// Package config - пакет формирования конфигурации.
package config

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name     string
		filepath string
		want     *Config
		wantErr  bool
	}{
		{
			name:     "Negative: проверяем наличие конфигурационного файла",
			filepath: "",
			want:     nil,
			wantErr:  true,
		},
		{
			name:     "Positive: проверяем содержимое конфигурационного файла",
			filepath: "../../cmd/config.yml",
			want: &Config{
				HTTP: HTTP{
					Host:              "0.0.0.0:8080",
					ReadTimeout:       10 * time.Second,
					WriteTimeout:      15 * time.Second,
					ReadHeaderTimeout: 5 * time.Second,
				},
				App: App{
					CountMapItems: 1024,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			configPath = tt.filepath
			once = sync.Once{}

			got, err := New()

			assert.EqualValuesf(t, err != nil, tt.wantErr, "%v", err)
			assert.EqualValues(t, got, tt.want)

		})
	}
}
