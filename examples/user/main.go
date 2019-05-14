package main

import (
	"fmt"
	"github.com/chiguirez/eventbus"
	"github.com/pkg/errors"
)

type UserRegistered struct {
	email string
}

type PrintUserOnRegisteredUserEventSubscriber struct{}

func (h PrintUserOnRegisteredUserEventSubscriber) Handle(event eventbus.Event) {
	userRegistered, ok := event.(UserRegistered)
	if !ok {
		return
	}

	fmt.Println("registered", userRegistered.email)
}

type emailSender func(email string) error

type SendWelcomeEmailUserOnRegisteredUserEventSubscriber struct{
	emailSender emailSender
}

func (h SendWelcomeEmailUserOnRegisteredUserEventSubscriber) Handle(event eventbus.Event) {
	userRegistered, ok := event.(UserRegistered)
	if !ok {
		return
	}

	h.emailSender(userRegistered.email)
}

type LoggingMiddleware struct{}

func (m LoggingMiddleware) Execute(event eventbus.Event, next eventbus.EventCallable) error {
	fmt.Println("Execution of logging middleware")

	return next(event)
}

func main() {
	mapSubscriberResolver := eventbus.NewMapSubscriberResolver()
	mapSubscriberResolver.AddSubscriber(new(UserRegistered), new(PrintUserOnRegisteredUserEventSubscriber))
	emailSender := func(email string) error {
		fmt.Println("sending welcome email to", email)
		return nil
	}
	mapSubscriberResolver.AddSubscriber(new(UserRegistered), SendWelcomeEmailUserOnRegisteredUserEventSubscriber{emailSender})
	bus := eventbus.NewBus(&mapSubscriberResolver, new(LoggingMiddleware))
	event := UserRegistered{"some@email.com"}
	err := bus.Publish(event)
	if err != nil {
		fmt.Printf("WARNING! an error occurred '%s'", err.Error())
	}
}
