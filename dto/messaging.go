package dto

import "sync"

type MessageBrokerConsumer struct {
	MessageRouting string
	Callback       func(wg *sync.WaitGroup, msg CreateEvent) error
}
