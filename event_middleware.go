package eventbus

type EventCallable func(event Event) error

type Middleware interface {
	Execute(event Event, next EventCallable) error
}

type eventHandlingMiddleware struct {
	subscriberResolver EventSubscriberResolver
}

func (m eventHandlingMiddleware) Execute(event Event, next EventCallable) error{
	subscribers, err := m.subscriberResolver.Resolve(event)
	if err != nil {
		return err
	}

	// TODO go routines
	for _, subscriber := range subscribers {
		if err = subscriber.Handle(event); err!=nil{
			return err
		}
	}

	return next(event)
}

type MiddlewareList []Middleware

func NewMiddlewareList(eventSubscriber eventHandlingMiddleware) MiddlewareList {
	return []Middleware{eventSubscriber}
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
