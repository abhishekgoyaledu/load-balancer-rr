package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"go.uber.org/zap"

	"github.com/coda-payments/game_app/pkg/utils"
)

const (
	confPath = "GAME_APP_CONF_PATH"
)

var Logger *zap.Logger

var RequestCount int

type Server struct {
	Port         int `json:"port"`
	ReadTimeout  int `json:"readTimeout"`
	WriteTimeout int `json:"writeTimeout"`
}

var Config = &struct {
	Server                   Server `json:"server"`
	ThresholdForSlowResponse []int  `json:"thresholdForSlowResponse"`
	MockSlowResponse         bool   `json:"mockSlowResponse"`
}{}

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

	configFile, err := os.ReadFile(configPath)
	if err != nil {
		// Push alert here
		Logger.Fatal("Failed to read config file", zap.String("path", configPath), zap.Error(err))
	}

	if err := json.Unmarshal(configFile, &Config); err != nil {
		// Push alert here
		Logger.Fatal("Failed to unmarshal config file", zap.Error(err))
	}

	if err := utils.ValidatePort(Config.Server.Port); err != nil {
		// Push alert here
		Logger.Fatal("Invalid port for load balancer in config file", zap.Error(err))
	}
}

// GetPort return port if user provided in args
func GetPort() int {
	// os.Args contains the command-line arguments
	// os.Args[0] is the program name
	var (
		port int
		err  error
	)
	args := os.Args

	// Check if there are any arguments passed
	if len(args) < 2 {
		fmt.Println("No arguments provided. Please provide some arguments.")
		port = Config.Server.Port
	} else {
		// Convert string to int
		port, err = strconv.Atoi(args[1])
		if err != nil {
			//Logger.Fatal("Error converting string to int:", zap.Error(err))
			return -1
		}
	}
	return port
}
