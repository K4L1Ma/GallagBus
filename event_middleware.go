package eventbus

type EventCallable func(event Event) error

type Middleware interface {
	Execute(event Event, next EventCallable) error
}

type eventHandlingMiddleware struct {
	handlerResolver EventHandlerResolver
}

func (m eventHandlingMiddleware) Execute(event Event, next EventCallable) error{
	handlers, err := m.handlerResolver.Resolve(event)
	if err != nil {
		return err
	}

	// TODO go routines
	for _, handler := range handlers {
		if err = handler.Handle(event); err!=nil{
			return err
		}
	}

	return next(event)
}

type MiddlewareList []Middleware

func NewMiddlewareList(eventHandler eventHandlingMiddleware) MiddlewareList {
	return []Middleware{eventHandler}
}

func (m MiddlewareList) Queue(middleware ...Middleware) MiddlewareList {
	return append(m, middleware...)
}

func (m MiddlewareList) start(event Event) error {
	return m.getCallable(0)(event)
}

func (m MiddlewareList) lastIndex() int {
	return len(m) - 1
}

func (m MiddlewareList) getCallable(index int) EventCallable {
	lastCallable := func(event Event) error {return nil}
	if index > m.lastIndex() {
		return lastCallable
	}

	return func(event Event) error{
		middleware := m[index]

		return middleware.Execute(event, m.getCallable(index+1))
	}
}
