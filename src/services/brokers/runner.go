package brokers

import (
	"app/src/core/brokers"
	"context"
	"encoding/json"
	"errors"
	"log"
	"sync"
	"time"
)

type NotificationHandler func(ctx context.Context, notification brokers.Notification) error

type Subscription struct {
	topic    string
	group    string
	handlers map[string]NotificationHandler
}

func (s *Subscription) On(notificationType string, handler NotificationHandler) *Subscription {
	s.handlers[notificationType] = handler
	return s
}

func (s *Subscription) dispatch(ctx context.Context, message brokers.Message) error {
	var notification brokers.Notification
	if err := json.Unmarshal(message.Value, &notification); err != nil {
		log.Printf("consumer %q: skipping malformed message: %v", s.topic, err)
		return nil
	}

	handler, ok := s.handlers[notification.Type]
	if !ok {
		log.Printf("consumer %q: skipping unknown command type %q", s.topic, notification.Type)
		return nil
	}

	return handler(ctx, notification)
}

type Runner struct {
	broker        brokers.Broker
	restartDelay  time.Duration
	subscriptions []*Subscription
}

func NewRunner(broker brokers.Broker, restartDelay time.Duration) *Runner {
	return &Runner{broker: broker, restartDelay: restartDelay}
}

func (r *Runner) Subscribe(topic string, group string) *Subscription {
	subscription := &Subscription{
		topic:    topic,
		group:    group,
		handlers: make(map[string]NotificationHandler),
	}
	r.subscriptions = append(r.subscriptions, subscription)
	return subscription
}

func (r *Runner) Run(ctx context.Context) {
	var wg sync.WaitGroup
	for _, subscription := range r.subscriptions {
		wg.Add(1)
		go func(s *Subscription) {
			defer wg.Done()
			r.runSubscription(ctx, s)
		}(subscription)
	}
	wg.Wait()
}

func (r *Runner) runSubscription(ctx context.Context, subscription *Subscription) {
	for {
		if ctx.Err() != nil {
			return
		}

		err := r.broker.Consumer().Consume(ctx, subscription.topic, subscription.group, subscription.dispatch)
		if err == nil || errors.Is(err, context.Canceled) || ctx.Err() != nil {
			return
		}

		log.Printf("consumer %q stopped: %v; restarting in %s", subscription.topic, err, r.restartDelay)
		select {
		case <-ctx.Done():
			return
		case <-time.After(r.restartDelay):
		}
	}
}
