package repository

import (
	"github.com/portless-io/shared-packages/dto"
	"github.com/streadway/amqp"
)

type MessageBrokerRepository interface {
	PublishMessage(messageType string, message interface{}) error
	Consume(consumer dto.MessageBrokerConsumer)
	SetNewRabbitMqChannel(ch *amqp.Channel)
	GetChannel() *amqp.Channel
}
