package event

import (
	"os"
	"sync"
)

type LogEvent struct {
	Message string
}

type LogChannel struct {
	channel chan LogEvent
	file    *os.File
	wg      sync.WaitGroup
}

func NewLogChannel(filePath string, bufferSize int) (*LogChannel, error) {
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	lc := &LogChannel{
		channel: make(chan LogEvent, bufferSize),
		file:    f,
	}

	lc.wg.Add(1)
	go lc.StartWorker()

	return lc, nil
}

func (lc *LogChannel) StartWorker() {
	defer lc.wg.Done()
	for logEvent := range lc.channel {
		lc.file.WriteString(logEvent.Message + "\n")
	}
}

func (lc *LogChannel) Send(message string) {
	lc.channel <- LogEvent{Message: message}
}

func (lc *LogChannel) Close() error {
	close(lc.channel)
	lc.wg.Wait()
	return lc.file.Close()
}
