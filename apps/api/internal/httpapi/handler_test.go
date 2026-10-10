package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mimustafa/research-workflow-platform/apps/api/internal/workflow"
)

type recordingPublisher struct {
	workflows []workflow.Workflow
}

type failingPublisher struct{}

func (p *recordingPublisher) Publish(_ context.Context, item workflow.Workflow) error {
	p.workflows = append(p.workflows, item)
	return nil
}

func (failingPublisher) Publish(context.Context, workflow.Workflow) error {
	return errors.New("broker unavailable")
}

func TestCreateAndListWorkflow(t *testing.T) {
	publisher := &recordingPublisher{}
	handler := NewHandler(slog.New(slog.NewTextHandler(bytes.NewBuffer(nil), nil)), publisher)

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
	if len(publisher.workflows) != 1 {
		t.Fatalf("published workflows = %d, want 1", len(publisher.workflows))
	}
}

func TestCreateWorkflowRejectsMalformedJSON(t *testing.T) {
	handler := NewHandler(slog.Default(), &recordingPublisher{})
	request := httptest.NewRequest(http.MethodPost, "/api/v1/workflows", bytes.NewBufferString(`{"brief":`))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusBadRequest)
	}
}

func TestCreateWorkflowDoesNotStoreWhenPublishFails(t *testing.T) {
	handler := NewHandler(slog.Default(), failingPublisher{})
	request := httptest.NewRequest(http.MethodPost, "/api/v1/workflows", bytes.NewBufferString(`{"brief":"Unpublished"}`))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusServiceUnavailable)
	}
	listResponse := httptest.NewRecorder()
	handler.ServeHTTP(listResponse, httptest.NewRequest(http.MethodGet, "/api/v1/workflows", nil))
	if listResponse.Body.String() != "[]\n" {
		t.Fatalf("stored workflows = %s, want empty list", listResponse.Body.String())
	}
}
