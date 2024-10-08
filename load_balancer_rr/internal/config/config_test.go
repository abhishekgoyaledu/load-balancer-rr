package config

import (
	"context"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"

	"go.uber.org/zap"

	"github.com/stretchr/testify/require"
)

func TestInit_InvalidConfigFile(t *testing.T) {
	// Setup environment variable
	os.Setenv("LOAD_BALANCER_RR_CONF_PATH", "invalid_config.json")
	defer os.Unsetenv("LOAD_BALANCER_RR_CONF_PATH")

	// Write invalid content to the temporary file
	err := ioutil.WriteFile("invalid_config.json", []byte("invalid content"), 0644)
	require.NoError(t, err)
	defer os.Remove("invalid_config.json")

	// Reinitialize the package level variables
	Logger = nil
	Config = &struct {
		Server                         Server  `json:"server"`
		Backend                        Backend `json:"backend"`
		HealthCheckTickerTimeInSeconds int64   `json:"healthCheckTickerTimeInSeconds"`
	}{}

	// Capture log output
	logger, err := zap.NewDevelopment()
	require.NoError(t, err)
	zap.ReplaceGlobals(logger)
}

func TestGracefulShutdownConfig(t *testing.T) {
	server := &http.Server{
		Addr: ":8080",
	}
	ctx, cancel := context.WithCancel(context.Background())

	go GracefulShutdownConfig(ctx, server)

	// Simulate a context cancel
	cancel()

	// Give some time for the shutdown process to start
	time.Sleep(1 * time.Second)
}
