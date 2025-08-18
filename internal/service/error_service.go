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

// Handle envia a mensagem de erro para o channel de logs
func (h *errorHandler) HandleError(ctx context.Context, err error, message string) {
	if err != nil {
		timestamp := time.Now().Format(time.RFC3339)
		logMsg := fmt.Sprintf("%s [ERROR] %s | details: %v", timestamp, message, err)

		// Envia de forma não bloqueante
		select {
		case <-ctx.Done():
			// Contexto cancelado, descarta log
			return
		case h.errorChan <- logMsg:
			// Log enviado
		default:
			// Channel cheio, descarta log ou poderia salvar em buffer
		}
	}
}
