package health

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"github.com/ivmerk/goshop/platform/pkg/logger"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(map[string]string{"status": "ok"}); err != nil {
		logger.Error(r.Context(), "failed to write health response", zap.Error(err))
	}
}
