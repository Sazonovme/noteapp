package api

import (
	"encoding/json"
	"net/http"
	"noteapp/pkg/logger"
	"time"
)

func apiError(w http.ResponseWriter, r *http.Request, statusCode int, err error) {
	w.WriteHeader(statusCode)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
	}
	logger.NewLog("api - apiError()", 5, err, "OUT - ERR "+time.Now().Format("02.01 15:04:05"), err.Error())
}
