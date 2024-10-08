package healtcheck

import (
	"math/rand"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/coda-payments/game_app/internal/config"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	if config.Config.MockSlowResponse {
		config.RequestCount++
		randomInt := rand.Intn(100000)
		config.Logger.Info("Request received healthcheck : ", zap.Int("request_id: ", randomInt), zap.Int("threshold: ", config.RequestCount))
		if config.RequestCount >= config.Config.ThresholdForSlowResponse[0] && config.RequestCount <= config.Config.ThresholdForSlowResponse[1] {
			time.Sleep(50 * time.Second)
		}
		config.Logger.Info("response healthcheck : %v", zap.Int("request_id: ", randomInt))
	}
	return
}
