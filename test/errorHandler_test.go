package test

import (
	"context"
	"corebanking/internal/utils"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockErrorHandler struct {
	Called  bool
	Ctx     context.Context
	Err     error
	Message string
}

func (m *MockErrorHandler) Handle(ctx context.Context, err error, message string) {
	m.Called = true
	m.Ctx = ctx
	m.Err = err
	m.Message = message
}

func TestHandleHTTPError(t *testing.T) {
	mockHandler := &MockErrorHandler{}
	recorder := httptest.NewRecorder()
	testErr := errors.New("some error")
	testMsg := "Test message"

	utils.HandleHTTPError(recorder, testErr, testMsg, mockHandler)

	if !mockHandler.Called {
		t.Errorf("expected ErrorHandler to be called")
	}
	if mockHandler.Err != testErr {
		t.Errorf("expected err %v, got %v", testErr, mockHandler.Err)
	}
	if mockHandler.Message != testMsg {
		t.Errorf("expected message %s, got %s", testMsg, mockHandler.Message)
	}

	if recorder.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, recorder.Code)
	}

	var resp map[string]string
	err := json.NewDecoder(recorder.Body).Decode(&resp)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp["error"] != testMsg {
		t.Errorf("expected error message %s, got %s", testMsg, resp["error"])
	}
	if resp["details"] != testErr.Error() {
		t.Errorf("expected details %s, got %s", testErr.Error(), resp["details"])
	}
}

func TestHandleHTTPError_NoErr(t *testing.T) {
	mockHandler := &MockErrorHandler{}
	recorder := httptest.NewRecorder()
	testMsg := "Another test"

	utils.HandleHTTPError(recorder, nil, testMsg, mockHandler)

	if !mockHandler.Called {
		t.Errorf("expected ErrorHandler to be called")
	}

	var resp map[string]string
	err := json.NewDecoder(recorder.Body).Decode(&resp)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp["details"] != "" {
		t.Errorf("expected empty details, got %s", resp["details"])
	}
}
