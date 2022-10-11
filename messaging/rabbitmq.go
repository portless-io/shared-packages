package messaging

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/portless-io/shared-packages/dto"
	"github.com/streadway/amqp"
)

type RabbitMqRepository struct {
	Ch  *amqp.Channel
	url string
}

func NewRabbitMqRepository(ch *amqp.Channel, url string) *RabbitMqRepository {
	rabbitmqRepository := &RabbitMqRepository{
		Ch:  ch,
		url: url,
	}

	return rabbitmqRepository
}

func (rm *RabbitMqRepository) SetNewRabbitMqChannel(ch *amqp.Channel) {
	rm.Ch = ch
}

func (rm *RabbitMqRepository) PublishMessage(topic string, message interface{}) error {
	messageInByte, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("rabbitmq publish: failed marshalling msg")
	}

	return rm.Ch.Publish(
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

func (rm *RabbitMqRepository) Consume(consumer dto.MessageBrokerConsumer) {
	q, err := rm.Ch.QueueDeclare(
		consumer.MessageRouting, // name
		false,                   // durable
		false,                   // delete when unused
		false,                   // exclusive
		false,                   // no-wait
		nil,                     // arguments
	)

	failOnError(err, "Failed to declare a queue")

	msgs, err := rm.Ch.Consume(
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

func (rm *RabbitMqRepository) Init() {
	go func(url string) {
		for {
			time.Sleep(15 * time.Second)
			<-rm.Ch.NotifyClose(make(chan *amqp.Error))

			rabbitMQConnection, err := amqp.Dial(url)
			if err != nil {
				log.Printf("RabbitMQ: failed re-connect to broker: %s", err.Error())
				continue
			}

			log.Println("re-connected to broker")
			defer rabbitMQConnection.Close()

			rabbitMQChannel, err := rabbitMQConnection.Channel()
			if err != nil {
				log.Printf("RabbitMQ: failed re-open channel %s", err.Error())
				continue
			}

			rm.SetNewRabbitMqChannel(rabbitMQChannel)
			break
		}
	}(rm.url)
}

func (rm *RabbitMqRepository) GetChannel() *amqp.Channel {
	return rm.Ch
}
