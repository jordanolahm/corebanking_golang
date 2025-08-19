package test

import (
	"bufio"
	"corebanking/internal/event"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLogChannel(t *testing.T) {
	tmpFile := filepath.Join(os.TempDir(), "log_test.txt")
	defer os.Remove(tmpFile)

	lc, err := event.NewLogChannel(tmpFile, 10)
	if err != nil {
		t.Fatalf("failed to create LogChannel: %v", err)
	}

	messages := []string{"first message", "second message", "third message"}

	for _, msg := range messages {
		lc.Send(msg)
	}

	if err := lc.Close(); err != nil {
		t.Fatalf("failed to close LogChannel: %v", err)
	}

	file, err := os.Open(tmpFile)
	if err != nil {
		t.Fatalf("failed to open log file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if len(lines) != len(messages) {
		t.Fatalf("expected %d lines, got %d", len(messages), len(lines))
	}

	for i, msg := range messages {
		if strings.TrimSpace(lines[i]) != msg {
			t.Errorf("expected line %d to be %q, got %q", i, msg, lines[i])
		}
	}
}

func TestLogChannel_BufferOverflow(t *testing.T) {
	tmpFile := filepath.Join(os.TempDir(), "log_test_overflow.txt")
	defer os.Remove(tmpFile)

	lc, err := event.NewLogChannel(tmpFile, 2)
	if err != nil {
		t.Fatalf("failed to create LogChannel: %v", err)
	}

	for i := 0; i < 5; i++ {
		lc.Send("msg")
	}

	if err := lc.Close(); err != nil {
		t.Fatalf("failed to close LogChannel: %v", err)
	}

	if _, err := os.Stat(tmpFile); err != nil {
		t.Errorf("log file not found: %v", err)
	}
}
