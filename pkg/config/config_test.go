package config

import (
	"os"
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func resetFlags() {
	pflag.CommandLine = pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)
}

func TestLoad_DefaultValues(t *testing.T) {
	resetFlags()

	originalEnv := map[string]string{}
	envVars := []string{"PORT", "POSTGRES_HOST", "POSTGRES_PORT", "POSTGRES_USER", "POSTGRES_PASSWORD", "POSTGRES_DB"}

	for _, key := range envVars {
		originalEnv[key] = os.Getenv(key)
		os.Unsetenv(key)
	}
	defer func() {
		for key, value := range originalEnv {
			if value != "" {
				os.Setenv(key, value)
			}
		}
	}()

	cfg, err := Load()
	require.NoError(t, err)
	assert.NotNil(t, cfg)

	assert.Equal(t, 8080, cfg.Port)
	assert.Equal(t, "localhost", cfg.DB.Host)
	assert.Equal(t, 5432, cfg.DB.Port)
	assert.Equal(t, "myuser", cfg.DB.User)
	assert.Equal(t, "mypassword", cfg.DB.Password)
	assert.Equal(t, "test_db", cfg.DB.DBName)
}

func TestLoad_WithEnvironmentVariables(t *testing.T) {
	resetFlags()

	originalEnv := map[string]string{}
	envVars := []string{"PORT", "POSTGRES_HOST", "POSTGRES_PORT", "POSTGRES_USER", "POSTGRES_PASSWORD", "POSTGRES_DB"}

	for _, key := range envVars {
		originalEnv[key] = os.Getenv(key)
	}
	defer func() {
		for key, value := range originalEnv {
			if value != "" {
				os.Setenv(key, value)
			} else {
				os.Unsetenv(key)
			}
		}
	}()

	os.Setenv("PORT", "9090")
	os.Setenv("POSTGRES_HOST", "test-host")
	os.Setenv("POSTGRES_PORT", "5433")
	os.Setenv("POSTGRES_USER", "test-user")
	os.Setenv("POSTGRES_PASSWORD", "test-password")
	os.Setenv("POSTGRES_DB", "test-db")

	cfg, err := Load()
	require.NoError(t, err)
	assert.NotNil(t, cfg)

	assert.Equal(t, 9090, cfg.Port)
	assert.Equal(t, "test-host", cfg.DB.Host)
	assert.Equal(t, 5433, cfg.DB.Port)
	assert.Equal(t, "test-user", cfg.DB.User)
	assert.Equal(t, "test-password", cfg.DB.Password)
	assert.Equal(t, "test-db", cfg.DB.DBName)
}

func TestGetEnv(t *testing.T) {
	os.Setenv("TEST_KEY", "test-value")
	defer os.Unsetenv("TEST_KEY")

	value := os.Getenv("TEST_KEY")
	assert.Equal(t, "test-value", value)
}
