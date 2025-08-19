package utils

import (
	"context"
	"encoding/json"
	"net/http"
)

type ErrorHandler interface {
	Handle(ctx context.Context, err error, message string)
}

func HandleHTTPError(w http.ResponseWriter, err error, message string, logger ErrorHandler) {
	if logger != nil {
		logger.Handle(context.Background(), err, message)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)

	resp := map[string]string{
		"error":   message,
		"details": "",
	}
	if err != nil {
		resp["details"] = err.Error()
	}
	json.NewEncoder(w).Encode(resp)
}
