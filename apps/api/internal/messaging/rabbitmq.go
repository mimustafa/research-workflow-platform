package messaging

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/mimustafa/research-workflow-platform/apps/api/internal/workflow"
)

const QueueName = "workflow.jobs"

type RabbitPublisher struct {
	channel *amqp.Channel
}

func NewRabbitPublisher(channel *amqp.Channel) (*RabbitPublisher, error) {
	if _, err := channel.QueueDeclare(QueueName, true, false, false, false, nil); err != nil {
		return nil, fmt.Errorf("declare queue: %w", err)
	}
	return &RabbitPublisher{channel: channel}, nil
}

func (p *RabbitPublisher) Publish(ctx context.Context, item workflow.Workflow) error {
	body, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("encode workflow job: %w", err)
	}

	return p.channel.PublishWithContext(ctx, "", QueueName, false, false, amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		Body:         body,
	})
}
