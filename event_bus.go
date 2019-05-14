package eventbus

type Event interface{}

type EventHandler interface {
	Handle(event Event) error
}

type EventBus interface {
	Publish(event Event) error
}

type Bus struct {
	middlewares MiddlewareList
}

func NewBus(handlerResolver EventHandlerResolver, middlewares ...Middleware) Bus {
	eventHandlingMiddleware := eventHandlingMiddleware{handlerResolver}
	middlewareList := NewMiddlewareList(eventHandlingMiddleware).Queue(middlewares...)

	return Bus{middlewareList}
}

func (b Bus) Publish(event Event) error {
	err := b.middlewares.start(event)

	return err
}
