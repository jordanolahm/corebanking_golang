package service

import (
	"context"
	"fmt"
	"time"
)

type ErrorHandler interface {
	HandleError(ctx context.Context, err error, message string)
}

type errorHandler struct {
	errorChan chan<- string
}

func NewErrorHandler(errorChan chan<- string) ErrorHandler {
	return &errorHandler{errorChan: errorChan}
}

func (h *errorHandler) HandleError(ctx context.Context, err error, message string) {
	if err != nil {
		timestamp := time.Now().Format(time.RFC3339)
		logMsg := fmt.Sprintf("%s [ERROR] %s | details: %v", timestamp, message, err)

		select {
		case <-ctx.Done():
			return
		case h.errorChan <- logMsg:
		default:
		}
	}
}
