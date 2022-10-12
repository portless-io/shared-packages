package messaging

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/portless-io/shared-packages/dto"
	"github.com/portless-io/shared-packages/repository"
	"github.com/streadway/amqp"
)

type rabbitMqRepository struct {
	ch        *amqp.Channel
	url       string
	consumers *[]dto.MessageBrokerConsumer
}

func NewRabbitMqRepository(url string, consumers *[]dto.MessageBrokerConsumer) repository.MessageBrokerRepository {
	rabbitMQConnection, err := amqp.Dial(url)
	if err != nil {
		log.Fatalf("RabbitMQ: failed connect to broker: %s", err.Error())
		panic(err)
	}

	log.Println("connected to broker")

	rabbitMQChannel, err := rabbitMQConnection.Channel()
	if err != nil {
		log.Fatalf("RabbitMQ: failed open channel")
		panic(err)
	}

	messagingRepository := &rabbitMqRepository{
		ch:        rabbitMQChannel,
		url:       url,
		consumers: consumers,
	}

	if consumers != nil {
		for _, consumer := range *consumers {
			messagingRepository.Consume(consumer)
		}
	}

	go func(url string) {
		for {
			time.Sleep(20 * time.Second)
			<-rabbitMQChannel.NotifyClose(make(chan *amqp.Error))

			log.Println("trying to re-connect to message broker")
			rabbitMQConnection, err := amqp.Dial(url)
			if err != nil {
				log.Printf("RabbitMQ: failed re-connect to broker: %s", err.Error())
				continue
			}

			log.Println("re-connected to message broker")

			rabbitMQChannel, err := rabbitMQConnection.Channel()
			if err != nil {
				log.Printf("RabbitMQ: failed re-open channel %s", err.Error())
				continue
			}

			messagingRepository.SetNewRabbitMqChannel(rabbitMQChannel)

			if messagingRepository.consumers != nil {
				for _, consumer := range *messagingRepository.consumers {
					messagingRepository.Consume(consumer)
				}
			}
		}
	}(url)

	return messagingRepository
}

func (rm *rabbitMqRepository) SetNewRabbitMqChannel(ch *amqp.Channel) {
	rm.ch = ch
}

func (rm *rabbitMqRepository) PublishMessage(topic string, message interface{}) error {
	messageInByte, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("rabbitmq publish: failed marshalling msg")
	}

	return rm.ch.Publish(
		"",    // exchange
		topic, // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        messageInByte,
		})
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func (rm *rabbitMqRepository) Consume(consumer dto.MessageBrokerConsumer) {
	q, err := rm.ch.QueueDeclare(
		consumer.MessageRouting, // name
		false,                   // durable
		false,                   // delete when unused
		false,                   // exclusive
		false,                   // no-wait
		nil,                     // arguments
	)

	failOnError(err, "Failed to declare a queue")

	msgs, err := rm.ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)

	failOnError(err, "Failed to register a consumer")

	maxNbConcurrentGoroutines := 5
	concurrentGoroutines := make(chan struct{}, maxNbConcurrentGoroutines)

	var wg sync.WaitGroup
	go func() {
		for d := range msgs {
			var data dto.CreateEvent

			err := json.Unmarshal(d.Body, &data)
			if err != nil {
				log.Println("consumer failed: ", err.Error())
				break
			}

			go func() {
				wg.Add(1)

				concurrentGoroutines <- struct{}{}
				err = consumer.Callback(&wg, data)
				if err != nil {
					log.Println("consumer failed: ", err.Error())
				}
				<-concurrentGoroutines
			}()
		}
	}()
	wg.Wait()
}

func (rm *rabbitMqRepository) GetChannel() *amqp.Channel {
	return rm.ch
}

func (rm *rabbitMqRepository) AddConsumers(consumers []dto.MessageBrokerConsumer) {
	if consumers != nil {
		for _, consumer := range consumers {
			rm.Consume(consumer)
		}
	}

	rm.consumers = &consumers
}

func (rm *rabbitMqRepository) CloseChannel() error {
	return rm.ch.Close()
}
