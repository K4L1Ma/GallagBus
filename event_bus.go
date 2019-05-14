package eventbus

type Event interface{}

type EventSubscriber interface {
	Handle(event Event)
}

type EventBus interface {
	Publish(event Event) error
}

type Bus struct {
	middlewares MiddlewareList
}

func NewBus(subscriberResolver EventSubscriberResolver, middlewares ...Middleware) Bus {
	eventHandlingMiddleware := eventHandlingMiddleware{subscriberResolver}
	middlewareList := NewMiddlewareList(eventHandlingMiddleware).Queue(middlewares...)

	return Bus{middlewareList}
}

func (b Bus) Publish(event Event) error {
	err := b.middlewares.start(event)

	return err
}
