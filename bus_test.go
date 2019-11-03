package gallagbus

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestGallagBus(t *testing.T) {
	r := require.New(t)

	CallChannel := make(chan struct{})
	var CallNumber uint64

	const eventType = "test"
	t.Run("Given a GallagBus and a EventListener Subscribed using a regex", func(t *testing.T) {
		eventBus := New()
		listener := NewEventListener(func(struct{}) {
			CallChannel <- struct{}{}
		}, QueueSize(1))
		eventBus.Subscribe("t[a-z]+t", listener)
		t.Run("When we publish an event with eventType test", func(t *testing.T) {
			eventBus.Publish(eventType, struct{}{})
			t.Run("Then listener got called 1 Time", func(t *testing.T) {
				for {
					select {
					case <-time.After(time.Millisecond * 1):
						t.Fatal("listener Should have been Called once")
						return
					case i := <-CallChannel:
						r.Equal(i, struct{}{}, "listener Should have been Called once")
						return
					}
				}
			})
		})
	})

	t.Run("Given a GallagBus and a 4 EventListener Subscribed into the BUS to the eventType test", func(t *testing.T) {
		eventBus := New()
		listener := NewEventListener(func(struct{}) {
			atomic.AddUint64(&CallNumber, 1)
			if atomic.LoadUint64(&CallNumber) >= 4 {
				CallChannel <- struct{}{}
			}
		}, QueueSize(1))
		eventBus.Subscribe(eventType, listener)
		listener1 := NewEventListener(func(struct{}) {
			atomic.AddUint64(&CallNumber, 1)
			if atomic.LoadUint64(&CallNumber) >= 4 {
				CallChannel <- struct{}{}
			}
		})
		eventBus.Subscribe(eventType, listener1)
		listener2 := NewEventListener(func(struct{}) {
			atomic.AddUint64(&CallNumber, 1)
			if atomic.LoadUint64(&CallNumber) >= 4 {
				CallChannel <- struct{}{}
			}
		})
		eventBus.Subscribe(eventType, listener2)
		listener3 := NewEventListener(func(struct{}) {
			atomic.AddUint64(&CallNumber, 1)
			if atomic.LoadUint64(&CallNumber) >= 4 {
				CallChannel <- struct{}{}
			}
		})
		eventBus.Subscribe(eventType, listener3)
		t.Run("When we publish an event with eventType test", func(t *testing.T) {
			eventBus.Publish(eventType, struct{}{})
			t.Run("Then listener get 4 Calls", func(t *testing.T) {
				<-CallChannel
				callNumber := atomic.LoadUint64(&CallNumber)
				r.Equal(callNumber, uint64(4), "Listener Should have been Called once")
			})
		})
	})

	t.Run("Given a GallagBus and a EventListener Subscribed into the BUS to the eventType test", func(t *testing.T) {
		eventBus := New()
		listener := NewEventListener(func(struct{}) {
			CallChannel <- struct{}{}
		}, QueueSize(1))
		eventBus.Subscribe(eventType, listener)
		t.Run("When we publish an other eventType", func(t *testing.T) {
			otherType := "other"
			eventBus.Publish(otherType, struct{}{})
			t.Run("Then after a few seconds listener never get Called", func(t *testing.T) {
				for {
					select {
					case <-time.After(10 * time.Millisecond):
						return
					case <-CallChannel:
						t.Fatal("listener shouldn't have been called")
					}
				}
			})
		})
	})

	t.Run("Given a GallagBus and a EventListener Subscribed into the BUS to the eventType test", func(t *testing.T) {
		eventBus := New()
		eventType := "test"
		listener := NewEventListener(func(struct{}) {
			CallChannel <- struct{}{}
		}, QueueSize(1))
		eventBus.Subscribe(eventType, listener)
		t.Run("When we shutdown graceful", func(t *testing.T) {
			eventBus.GracefulShutdown()
			t.Run("Then publishing should cause a panic", func(t *testing.T) {
				r.Panics(func() { eventBus.Publish(eventType, struct{}{}) })
			})
		})
	})
}
