package eventbus

import (
	"fmt"
	"github.com/chiguirez/eventbus/typer"
)

//go:generate moq -out event_subscriber_resolver_mock.go . EventSubscriberResolver
type EventSubscriberResolver interface {
	Resolve(event Event) ([]EventSubscriber, error)
}

type MapSubscriberResolver struct {
	subscribers map[string][]EventSubscriber
}

func NewMapSubscriberResolver() MapSubscriberResolver {
	return MapSubscriberResolver{
		map[string][]EventSubscriber{},
	}
}

func (r MapSubscriberResolver) Resolve(event Event) ([]EventSubscriber, error) {
	subscriber, ok := r.subscribers[typer.Identify(event)]
	if !ok {
		return nil, fmt.Errorf("could not find event subscribers")
	}

	return subscriber, nil
}

func (r MapSubscriberResolver) AddSubscriber(event Event, subscriber EventSubscriber) {
	subscribers := r.subscribers[typer.Identify(event)]
	if len(subscribers) == 0 {
		r.subscribers[typer.Identify(event)] = []EventSubscriber{subscriber}
		return
	}

	r.subscribers[typer.Identify(event)] = append(subscribers, subscriber)
}
