package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mimustafa/research-workflow-platform/apps/api/internal/httpapi"
	"github.com/mimustafa/research-workflow-platform/apps/api/internal/messaging"
	"github.com/mimustafa/research-workflow-platform/apps/api/internal/store"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://research:research_local_password@localhost:5432/research_workflows"
	}
	pool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		logger.Error("configure database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()
	if err := pool.Ping(context.Background()); err != nil {
		logger.Error("connect to database", "error", err)
		os.Exit(1)
	}

	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://research:research_local_password@localhost:5672/"
	}

	connection, err := amqp.Dial(rabbitURL)
	if err != nil {
		logger.Error("connect to rabbitmq", "error", err)
		os.Exit(1)
	}
	defer connection.Close()

	channel, err := connection.Channel()
	if err != nil {
		logger.Error("open rabbitmq channel", "error", err)
		os.Exit(1)
	}
	defer channel.Close()

	publisher, err := messaging.NewRabbitPublisher(channel)
	if err != nil {
		logger.Error("configure rabbitmq publisher", "error", err)
		os.Exit(1)
	}

	handler := httpapi.NewHandler(logger, publisher, store.NewPostgres(pool))

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
