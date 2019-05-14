package eventbus

import (
	"github.com/stretchr/testify/require"
	"testing"
)

type AEventHandler struct {
	NumberOfHandleCalls int
}

func (h *AEventHandler) Handle(event Event) error {
	h.NumberOfHandleCalls++
	return nil
}

type AMiddleware struct {
	NumberOfExecuteCalls int
}

func (m *AMiddleware) Execute(event Event, next EventCallable) error {
	m.NumberOfExecuteCalls++

	return next(event)
}

func TestMapHandlerResolver_Resolve(t *testing.T) {
	isRequire := require.New(t)
	event := struct{}{}
	t.Run("Given a map handler resolver", func(t *testing.T) {
		sut := NewMapHandlerResolver()
		t.Run("When event handler is not found", func(t *testing.T) {
			handler, err := sut.Resolve(event)
			t.Run("Then an error is returned", func(t *testing.T) {
				isRequire.Nil(handler)
				isRequire.Error(err)
			})
		})
		t.Run("When a event with its handler is added", func(t *testing.T) {
			handler := &AEventHandler{}
			sut.AddHandler(event, handler)
			t.Run("Then the handler is resolved", func(t *testing.T) {
				resolvedHandlers, err := sut.Resolve(event)
				isRequire.Equal(handler, resolvedHandlers[0])
				isRequire.NoError(err)
			})
		})
	})
}

func TestEventbus_Dispatch(t *testing.T) {
	isRequire := require.New(t)
	t.Run("Given a eventbus event bus without middlewares", func(t *testing.T) {
		handler := AEventHandler{}
		handlerResolver := EventHandlerResolverMock{
			ResolveFunc: func(event Event) ([]EventHandler, error) {
				return []EventHandler{&handler}, nil
			},
		}
		sut := NewBus(&handlerResolver)
		t.Run("When a event is dispatched", func(t *testing.T) {
			event := struct{}{}
			sut.Publish(event)
			t.Run("Then the resolved event handler handles the event", func(t *testing.T) {
				resolverHasBeenCalled := len(handlerResolver.ResolveCalls()) > 0
				isRequire.True(resolverHasBeenCalled)
				handlerHasBeenCalled := handler.NumberOfHandleCalls > 0
				isRequire.True(handlerHasBeenCalled)
			})
		})
	})
	t.Run("Given a eventbus event bus with middlewares", func(t *testing.T) {
		aMiddleware := &AMiddleware{}
		anotherMiddleware := &AMiddleware{}
		handler := AEventHandler{}
		handlerResolver := EventHandlerResolverMock{
			ResolveFunc: func(event Event) ([]EventHandler, error) {
				return []EventHandler{&handler}, nil
			},
		}
		sut := NewBus(&handlerResolver, aMiddleware, anotherMiddleware)
		t.Run("When a event is dispatched", func(t *testing.T) {
			event := struct{}{}
			sut.Publish(event)
			t.Run("Then the resolved event handler handles the event", func(t *testing.T) {
				resolverHasBeenCalled := len(handlerResolver.ResolveCalls()) > 0
				isRequire.True(resolverHasBeenCalled)
				handlerHasBeenCalled := handler.NumberOfHandleCalls > 0
				isRequire.True(handlerHasBeenCalled)
			})
			t.Run("And the middlewares are executed", func(t *testing.T) {
				isRequire.True(aMiddleware.NumberOfExecuteCalls > 0)
				isRequire.True(anotherMiddleware.NumberOfExecuteCalls > 0)
			})
		})
	})
}
