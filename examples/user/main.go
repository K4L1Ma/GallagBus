package main

import (
	"fmt"
	"github.com/chiguirez/eventbus"
	"github.com/pkg/errors"
)

type UserRegistered struct {
	email string
}

type PrintUserOnRegisteredUserEventHandler struct{}

func (h PrintUserOnRegisteredUserEventHandler) Handle(event eventbus.Event) error {
	userRegistered, ok := event.(UserRegistered)
	if !ok {
		return errors.New("Could not handle a non register user event")
	}

	fmt.Println("registered", userRegistered.email)

	return nil
}

type emailSender func(email string) error

type SendWelcomeEmailUserOnRegisteredUserEventHandler struct{
	emailSender emailSender
}

func (h SendWelcomeEmailUserOnRegisteredUserEventHandler) Handle(event eventbus.Event) error {
	userRegistered, ok := event.(UserRegistered)
	if !ok {
		return errors.New("Could not handle a non register user event")
	}

	return h.emailSender(userRegistered.email)
}

type LoggingMiddleware struct{}

func (m LoggingMiddleware) Execute(event eventbus.Event, next eventbus.EventCallable) error {
	fmt.Println("Execution of logging middleware")

	return next(event)
}

func main() {
	mapHandlerResolver := eventbus.NewMapHandlerResolver()
	mapHandlerResolver.AddHandler(new(UserRegistered), new(PrintUserOnRegisteredUserEventHandler))
	emailSender := func(email string) error {
		fmt.Println("sending welcome email to", email)
		return nil
	}
	mapHandlerResolver.AddHandler(new(UserRegistered), SendWelcomeEmailUserOnRegisteredUserEventHandler{emailSender})
	bus := eventbus.NewBus(&mapHandlerResolver, new(LoggingMiddleware))
	event := UserRegistered{"some@email.com"}
	err := bus.Publish(event)
	if err != nil {
		fmt.Printf("WARNING! an error occurred '%s'", err.Error())
	}
}
