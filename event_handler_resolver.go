package eventbus

import (
	"fmt"
	"github.com/chiguirez/eventbus/typer"
)

//go:generate moq -out event_handler_resolver_mock.go . EventHandlerResolver
type EventHandlerResolver interface {
	Resolve(event Event) ([]EventHandler, error)
}

type MapHandlerResolver struct {
	handlers map[string][]EventHandler
}

func NewMapHandlerResolver() MapHandlerResolver {
	return MapHandlerResolver{
		map[string][]EventHandler{},
	}
}

func (r MapHandlerResolver) Resolve(event Event) ([]EventHandler, error) {
	handler, ok := r.handlers[typer.Identify(event)]
	if !ok {
		return nil, fmt.Errorf("could not find event handlers")
	}

	return handler, nil
}

func (r MapHandlerResolver) AddHandler(event Event, handler EventHandler) {
	handlers := r.handlers[typer.Identify(event)]
	if len(handlers) == 0 {
		r.handlers[typer.Identify(event)] = []EventHandler{handler}
		return
	}

	r.handlers[typer.Identify(event)] = append(handlers, handler)
}
