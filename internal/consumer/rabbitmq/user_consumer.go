package rabbitmq

import (
	"bytes"
	"context"
	"encoding/gob"
	"os"
	"os/signal"
	"syscall"

	amqp "github.com/rabbitmq/amqp091-go"

	"post/internal/consumer"
	"post/internal/model"
	"post/internal/service"
	"post/pkg/logger"
)

const (
	queueName    = "user"
	ExchangeName = "user_events"
)

type Consumer struct {
	ch        *amqp.Channel
	processor service.EventsService
}

func NewUserConsumer(ch *amqp.Channel, processor service.EventsService, serviceName string) consumer.Consumer {

	queue, err := ch.QueueDeclare(
		serviceName,
		true,  // durable
		false, // auto delete
		false, // exclusive
		false, // no wait
		nil,   // args
	)
	if err != nil {
		return nil
	}

	// Привязываем очередь к exchange
	err = ch.QueueBind(
		queue.Name,
		"", // routing key не используется для fanout
		ExchangeName,
		false,
		nil,
	)
	if err != nil {
		return nil
	}

	return &Consumer{
		ch:        ch,
		processor: processor,
	}
}

func (u Consumer) Start(ctx context.Context, queueName string) {
	const op = "rabbitmq.Start"
	messages, err := u.ch.Consume(queueName, "", true, false, false, false, nil)
	if err != nil {
		logger.Fatal(op, "err: %s\n", err)
	}

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for {
			select {
			case <-ctx.Done():
				logger.Info("stopping event consumer by ctx")
				return
			case message := <-messages:
				logger.Info("Message: %s\n", message.Body)
				e := event(message.Body)
				err = u.processor.Process(ctx, *e)
				if err != nil {
					logger.Warn(op, "err: %s\n", err)
					logger.Warn(op, "event: %s\n", e)
				}
			case <-sigchan:
				logger.Info("Interrupt detected!")
				os.Exit(0)
			}
		}
	}()
}

func event(msg []byte) *model.Event {
	var e model.Event
	buf := bytes.NewBuffer(msg)

	if err := gob.NewDecoder(buf).Decode(&e); err != nil {
		return nil
	}

	logger.Info("event:", "e", e)

	return &e
}
