package httpapi

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateAndListWorkflow(t *testing.T) {
	handler := NewHandler(slog.New(slog.NewTextHandler(bytes.NewBuffer(nil), nil)))

	createRequest := httptest.NewRequest(http.MethodPost, "/api/v1/workflows", bytes.NewBufferString(`{"brief":"Find related studies"}`))
	createRequest.Header.Set("Content-Type", "application/json")
	createResponse := httptest.NewRecorder()
	handler.ServeHTTP(createResponse, createRequest)

	if createResponse.Code != http.StatusCreated {
		t.Fatalf("create status = %d, want %d", createResponse.Code, http.StatusCreated)
	}

	listResponse := httptest.NewRecorder()
	handler.ServeHTTP(listResponse, httptest.NewRequest(http.MethodGet, "/api/v1/workflows", nil))

	if listResponse.Code != http.StatusOK {
		t.Fatalf("list status = %d, want %d", listResponse.Code, http.StatusOK)
	}

	var workflows []struct {
		Brief string `json:"brief"`
	}
	if err := json.NewDecoder(listResponse.Body).Decode(&workflows); err != nil {
		t.Fatalf("decode list response: %v", err)
	}
	if len(workflows) != 1 || workflows[0].Brief != "Find related studies" {
		t.Fatalf("workflows = %+v, want one submitted workflow", workflows)
	}
}

func TestCreateWorkflowRejectsMalformedJSON(t *testing.T) {
	handler := NewHandler(slog.Default())
	request := httptest.NewRequest(http.MethodPost, "/api/v1/workflows", bytes.NewBufferString(`{"brief":`))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusBadRequest)
	}
}
