package event

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/velmie/broker"
)

type CardCreated struct {
	CardID string `json:"cardId"`
	UserID string `json:"userId"`
}

type BrokerPublisher struct {
	publisher broker.Publisher
}

func NewBrokerPublisher(publisher broker.Publisher) *BrokerPublisher {
	return &BrokerPublisher{publisher: publisher}
}

func (b *BrokerPublisher) OnCardCreated(_ context.Context, topicName string, msg *CardCreated) error {
	event, err := json.Marshal(msg)
	if err != nil {
		return errors.Wrap(err, "cannot marshall event into json")
	}

	message := &broker.Message{
		ID:   uuid.New().String(),
		Body: event,
	}

	if err = b.publisher.Publish(topicName, message); err != nil {
		return errors.Wrap(err, "cannot publish message")
	}
	return nil
}
