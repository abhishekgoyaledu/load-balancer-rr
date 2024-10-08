package routes

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/coda-payments/game_app/handler/healtcheck"
	"github.com/coda-payments/game_app/handler/player"
	"github.com/coda-payments/game_app/internal/config"
	"github.com/gorilla/mux"
)

func Launch(port int) {

	v1Router := mux.NewRouter()
	healthcheckAPI(v1Router)
	routeAPIs(v1Router)
	server := &http.Server{
		Handler: v1Router,
		Addr:    fmt.Sprintf(":%d", port),
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: time.Duration(config.Config.Server.WriteTimeout) * time.Second,
		ReadTimeout:  time.Duration(config.Config.Server.ReadTimeout) * time.Second,
	}

	config.Logger.Info("server started on : ", zap.Int("Port : ", port))
	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		// Push alert here: Launch encountered an unexpected error
		config.Logger.Fatal("Launch encountered an unexpected error", zap.Error(err))
	}
}

func healthcheckAPI(router *mux.Router) {
	router.HandleFunc("/healthcheck", healtcheck.HealthCheck).Methods("GET")
}

func routeAPIs(router *mux.Router) {
	router.HandleFunc("/create", player.CreatePlayer).Methods("POST")
}
