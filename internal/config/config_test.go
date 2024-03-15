package config

import (
	"os"
	"reflect"
	"testing"
)

func TestNewConfig(t *testing.T) {
	t.Run("Successful configuration", func(t *testing.T) {
		os.Setenv("PORT", "8081")
		os.Setenv("GITHUB_BASE_URL", "https://api.github.com")
		os.Setenv("GITHUB_TOKEN", "my-token")

		expectedCfg := Config{
			Port: 8081,
			Github: GithubConfig{
				BaseURL: "https://api.github.com",
				Token:   "my-token",
			},
		}

		cfg, err := NewConfig()
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		if !reflect.DeepEqual(cfg, expectedCfg) {
			t.Errorf("Expected configuration: %+v, but got: %+v", expectedCfg, cfg)
		}
	})
}
