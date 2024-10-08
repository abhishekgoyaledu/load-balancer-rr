package player

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"github.com/coda-payments/game_app/internal/config"
)

func CreatePlayer(w http.ResponseWriter, r *http.Request) {
	//push an alert here for latency
	var payload map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	config.Logger.Info("Request received : ", zap.String("request : ", fmt.Sprintf("%v", payload)))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payload)
}
