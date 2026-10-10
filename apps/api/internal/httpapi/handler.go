package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/mimustafa/research-workflow-platform/apps/api/internal/workflow"
)

type Handler struct {
	logger    *slog.Logger
	publisher Publisher
	mu        sync.RWMutex
	workflows []workflow.Workflow
}

type Publisher interface {
	Publish(context.Context, workflow.Workflow) error
}

type createWorkflowRequest struct {
	Brief string `json:"brief"`
}

func NewHandler(logger *slog.Logger, publisher Publisher) http.Handler {
	handler := &Handler{logger: logger, publisher: publisher}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handler.health)
	mux.HandleFunc("GET /api/v1/workflows", handler.listWorkflows)
	mux.HandleFunc("POST /api/v1/workflows", handler.createWorkflow)
	return withCORS(mux)
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (h *Handler) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) listWorkflows(w http.ResponseWriter, _ *http.Request) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	workflows := append([]workflow.Workflow(nil), h.workflows...)
	if workflows == nil {
		workflows = []workflow.Workflow{}
	}
	writeJSON(w, http.StatusOK, workflows)
}

func (h *Handler) createWorkflow(w http.ResponseWriter, r *http.Request) {
	var request createWorkflowRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "request body must be valid JSON")
		return
	}

	newWorkflow, err := workflow.New(strings.TrimSpace(request.Brief), time.Now())
	if errors.Is(err, workflow.ErrEmptyBrief) {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err != nil {
		h.logger.Error("create workflow failed", "error", err)
		writeError(w, http.StatusInternalServerError, "could not create workflow")
		return
	}

	if err := h.publisher.Publish(r.Context(), newWorkflow); err != nil {
		h.logger.Error("publish workflow failed", "workflow_id", newWorkflow.ID, "error", err)
		writeError(w, http.StatusServiceUnavailable, "workflow queue is unavailable")
		return
	}

	h.mu.Lock()
	h.workflows = append(h.workflows, newWorkflow)
	h.mu.Unlock()

	writeJSON(w, http.StatusCreated, newWorkflow)
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
