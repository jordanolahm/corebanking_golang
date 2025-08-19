package worker

import (
	"context"
	"corebanking/internal/event"
	"fmt"
	"time"
)

type ErrorWorker struct {
	logChannel *event.LogChannel
}

func NewErrorWorker(logChannel *event.LogChannel) *ErrorWorker {
	return &ErrorWorker{logChannel: logChannel}
}

func (w *ErrorWorker) Handle(ctx context.Context, err error, message string) {

	var logMsg string
	if err != nil {
		logMsg = fmt.Sprintf("[ERROR] %s | details: %v | time: %s",
			message, err, time.Now().Format(time.RFC3339))
	} else {
		logMsg = fmt.Sprintf("[ERROR] %s | time: %s",
			message, time.Now().Format(time.RFC3339))
	}

	select {
	case <-ctx.Done():
		return
	default:
		w.logChannel.Send(logMsg)
	}
}
