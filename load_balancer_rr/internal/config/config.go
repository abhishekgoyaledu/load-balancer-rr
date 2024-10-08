package config

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"go.uber.org/zap"

	"github.com/coda-payments/load_balancer_rr/pkg/utils"
)

const (
	confPath = "LOAD_BALANCER_RR_CONF_PATH"
)

var Logger *zap.Logger

// Config holds the overall configuration for the application, including server settings, backend configurations, and health check intervals.
var Config = &struct {
	Server  Server  `json:"server"`
	Backend Backend `json:"backend"`

	// HealthCheckTickerTimeInSeconds defines the interval for health check ticks in seconds.
	HealthCheckTickerTimeInSeconds int64 `json:"healthCheckTickerTimeInSeconds"`
}{}

// Server represents the configuration for the server settings.
type Server struct {
	Port         int `json:"port"`
	ReadTimeout  int `json:"readTimeout"`
	WriteTimeout int `json:"writeTimeout"`
}

// Backend holds the configuration for backend services, including server router and endpoints.
type Backend struct {
	Routes   []string            `json:"routes"`
	Endpoint map[string]Endpoint `json:"endpoints"`
}

// Endpoint defines the configuration for a single backend endpoint.
type Endpoint struct {
	URL     string `json:"url"`
	Timeout int    `json:"timeout"`
}

func init() {
	var err error
	Logger, err = zap.NewProduction()
	if err != nil {
		// Push alert here
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	configPath := os.Getenv(confPath)
	if configPath == "" {
		// Push alert here
		Logger.Fatal("Config path not set in environment variable", zap.String("envVar", confPath))
	}

	configFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		// Push alert here
		Logger.Fatal("Failed to read config file", zap.String("path", configPath), zap.Error(err))
	}

	if err := json.Unmarshal(configFile, &Config); err != nil {
		// Push alert here
		Logger.Fatal("Failed to unmarshal config file", zap.Error(err))
	}

	// Validating the port -> filters out currently used port in system
	if err := utils.ValidatePort(Config.Server.Port); err != nil {
		// Push alert here
		Logger.Fatal("Invalid port for load balancer in config file", zap.Error(err))
	}
}

// GracefulShutdownConfig Shutdown server gracefully on context cancellation
func GracefulShutdownConfig(ctx context.Context, server *http.Server) {
	go func() {
		<-ctx.Done()

		//creating shutdown context with 20sec timeout
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		//context's resources are freed when the function exits
		defer cancel()
		if shutdownErr := server.Shutdown(shutdownCtx); shutdownErr != nil {
			Logger.Fatal("Gracefully Shutdown Interrupted : ", zap.String("shutdown err: ", shutdownErr.Error()))
		}
	}()
}
