package broker

import "post/internal/client/broker/rabbitmq"

type ClientMsgBroker interface {
	Connect() *rabbitmq.RabbitMQ
	Close() error
}
