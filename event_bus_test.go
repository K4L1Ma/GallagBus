package eventbus

import (
	"github.com/stretchr/testify/require"
	"testing"
)

type AEventSubscriber struct {
	NumberOfSubscribeCalls int
}

func (h *AEventSubscriber) Handle(event Event) {
	h.NumberOfSubscribeCalls++
}

type AMiddleware struct {
	NumberOfExecuteCalls int
}

func (m *AMiddleware) Execute(event Event, next EventCallable) error {
	m.NumberOfExecuteCalls++

	return next(event)
}

func TestMapSubscriberResolver_Resolve(t *testing.T) {
	isRequire := require.New(t)
	event := struct{}{}
	t.Run("Given a map subscriber resolver", func(t *testing.T) {
		sut := NewMapSubscriberResolver()
		t.Run("When event subscriber is not found", func(t *testing.T) {
			subscriber, err := sut.Resolve(event)
			t.Run("Then an error is returned", func(t *testing.T) {
				isRequire.Nil(subscriber)
				isRequire.Error(err)
			})
		})
		t.Run("When a event with its subscriber is added", func(t *testing.T) {
			subscriber := &AEventSubscriber{}
			sut.AddSubscriber(event, subscriber)
			t.Run("Then the subscriber is resolved", func(t *testing.T) {
				resolvedSubscribers, err := sut.Resolve(event)
				isRequire.Equal(subscriber, resolvedSubscribers[0])
				isRequire.NoError(err)
			})
		})
	})
}

func TestEventbus_Dispatch(t *testing.T) {
	isRequire := require.New(t)
	t.Run("Given a eventbus event bus without middlewares", func(t *testing.T) {
		subscriber := AEventSubscriber{}
		subscriberResolver := EventSubscriberResolverMock{
			ResolveFunc: func(event Event) ([]EventSubscriber, error) {
				return []EventSubscriber{&subscriber}, nil
			},
		}
		sut := NewBus(&subscriberResolver)
		t.Run("When a event is dispatched", func(t *testing.T) {
			event := struct{}{}
			sut.Publish(event)
			t.Run("Then the resolved event subscriber subscribes the event", func(t *testing.T) {
				resolverHasBeenCalled := len(subscriberResolver.ResolveCalls()) > 0
				isRequire.True(resolverHasBeenCalled)
				subscriberHasBeenCalled := subscriber.NumberOfSubscribeCalls > 0
				isRequire.True(subscriberHasBeenCalled)
			})
		})
	})
	t.Run("Given a eventbus event bus with middlewares", func(t *testing.T) {
		aMiddleware := &AMiddleware{}
		anotherMiddleware := &AMiddleware{}
		subscriber := AEventSubscriber{}
		subscriberResolver := EventSubscriberResolverMock{
			ResolveFunc: func(event Event) ([]EventSubscriber, error) {
				return []EventSubscriber{&subscriber}, nil
			},
		}
		sut := NewBus(&subscriberResolver, aMiddleware, anotherMiddleware)
		t.Run("When a event is dispatched", func(t *testing.T) {
			event := struct{}{}
			sut.Publish(event)
			t.Run("Then the resolved event subscriber subscribes the event", func(t *testing.T) {
				resolverHasBeenCalled := len(subscriberResolver.ResolveCalls()) > 0
				isRequire.True(resolverHasBeenCalled)
				subscriberHasBeenCalled := subscriber.NumberOfSubscribeCalls > 0
				isRequire.True(subscriberHasBeenCalled)
			})
			t.Run("And the middlewares are executed", func(t *testing.T) {
				isRequire.True(aMiddleware.NumberOfExecuteCalls > 0)
				isRequire.True(anotherMiddleware.NumberOfExecuteCalls > 0)
			})
		})
	})
}
