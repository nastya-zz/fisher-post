package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	"post/internal/client/broker/rabbitmq"
	"post/internal/model"
	"post/pkg/logger"
)

const (
	exchangeName = "post_events"
	exchangeType = "fanout"
)

type PostProducer struct {
	client *rabbitmq.RabbitMQ
}

type Producer interface {
	PublishEvent(ctx context.Context, event *model.Event) error
}

func NewPostProducer(client *rabbitmq.RabbitMQ) (Producer, error) {
	producer := &PostProducer{
		client: client,
	}

	if err := producer.declareExchange(); err != nil {
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	return producer, nil
}

func (p *PostProducer) declareExchange() error {
	return p.client.Channel.ExchangeDeclare(
		exchangeName,
		exchangeType,
		true,  // durable
		false, // auto-deleted
		false, // internal
		false, // no-wait
		nil,   // arguments
	)
}

func (p *PostProducer) PublishEvent(ctx context.Context, event *model.Event) error {
	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	msg := amqp.Publishing{
		ContentType:  "application/json",
		Body:         body,
		Timestamp:    time.Now(),
		MessageId:    fmt.Sprintf("%d", event.ID),
		Type:         event.Type,
		DeliveryMode: amqp.Persistent,
	}

	err = p.client.Channel.PublishWithContext(
		ctx,
		exchangeName,
		"",    // routing key
		false, // mandatory
		false, // immediate
		msg,
	)

	if err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	logger.Info("Event published successfully",
		"event_id", event.ID,
		"event_type", event.Type,
	)

	return nil
}
