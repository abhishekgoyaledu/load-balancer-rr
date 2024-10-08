package config

import (
	"encoding/json"
	"os"
	"testing"

	"go.uber.org/zap"

	"github.com/coda-payments/game_app/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func setup(t *testing.T) (teardown func()) {
	// Save original environment variable and restore after test
	originalEnv := os.Getenv(confPath)
	teardown = func() {
		os.Setenv(confPath, originalEnv)
	}

	configPath := os.Getenv(confPath)

	// Mock the ValidatePort function
	utils.ValidatePort = func(port int) error {
		return nil
	}

	var err error

	// Reinitialize the config package
	Logger, err = zap.NewProduction()
	assert.NoError(t, err, "Failed to initialize logger")
	Config = &struct {
		Server                   Server `json:"server"`
		ThresholdForSlowResponse []int  `json:"thresholdForSlowResponse"`
		MockSlowResponse         bool   `json:"mockSlowResponse"`
	}{}

	configFile, err := os.ReadFile(configPath)
	assert.NoError(t, err, "Failed to read config file")
	err = json.Unmarshal(configFile, &Config)
	assert.NoError(t, err, "Failed to unmarshal config file")

	return teardown
}

func TestConfigInitialization(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	// Check if the configuration values are set correctly
	assert.Equal(t, 9000, Config.Server.Port)
	assert.Equal(t, 15, Config.Server.ReadTimeout)
	assert.Equal(t, 15, Config.Server.WriteTimeout)
	assert.Equal(t, []int{12, 25}, Config.ThresholdForSlowResponse)
	assert.True(t, Config.MockSlowResponse)
}

func TestGetPort(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	// Mock os.Args for testing
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	os.Args = []string{"program", "9090"}
	port := GetPort()
	assert.Equal(t, 9090, port, "Expected port to be 9090")

	os.Args = []string{"program"}
	port = GetPort()
	assert.Equal(t, 9000, port, "Expected port to be 9000 (default from config)")
}

func TestGetPortWithInvalidArgument(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	// Mock os.Args for testing
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	os.Args = []string{"program", "invalid"}
	port := GetPort()
	assert.Equal(t, port, -1)
}
