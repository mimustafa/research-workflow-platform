package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/mimustafa/research-workflow-platform/apps/api/internal/httpapi"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	handler := httpapi.NewHandler(logger)

	server := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	logger.Info("starting api", "address", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("api stopped", "error", err)
		os.Exit(1)
	}
}
