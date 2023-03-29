package event

import (
	"github.com/velmie/broker"
)

type SQSSubscriber struct {
	logger        broker.Logger
	subscriber    broker.Subscriber
	subscriptions []broker.Subscription
	eventHandlers map[string]broker.Handler
}

func NewSQSSubscriber(
	logger broker.Logger,
	sqsSubscriber broker.Subscriber,
) *SQSSubscriber {
	return &SQSSubscriber{
		logger:        logger,
		subscriber:    sqsSubscriber,
		eventHandlers: make(map[string]broker.Handler),
	}
}

func (s *SQSSubscriber) AddSubscription(queueName string, handler broker.Handler) {
	s.eventHandlers[queueName] = handler
}

func (s *SQSSubscriber) SubscribeAll() error {
	for topic, handler := range s.eventHandlers {
		subscription, err := s.subscriber.Subscribe(
			topic,
			handler,
			broker.WithDefaultErrorHandler(s.subscriber, s.logger),
		)
		if err != nil {
			return err
		}
		s.subscriptions = append(s.subscriptions, subscription)
	}
	return nil
}

func (s *SQSSubscriber) UnsubscribeAll() error {
	var errs Errors
	for i := range s.subscriptions {
		err := s.subscriptions[i].Unsubscribe()
		if err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return errs
}
