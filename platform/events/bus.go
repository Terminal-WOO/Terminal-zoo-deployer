package events

import (
	"context"
	"sync"
)

// EventBus handles event publishing and subscription
type EventBus interface {
	// Publish publishes an event
	Publish(ctx context.Context, event *Event) error

	// Subscribe subscribes to events of a specific type
	Subscribe(ctx context.Context, eventType EventType, handler EventHandler) error

	// Unsubscribe unsubscribes from events
	Unsubscribe(eventType EventType, handler EventHandler) error
}

// EventHandler handles events
type EventHandler func(ctx context.Context, event *Event) error

// InMemoryEventBus is an in-memory implementation of EventBus
type InMemoryEventBus struct {
	mu        sync.RWMutex
	handlers  map[EventType][]EventHandler
	events    []*Event
	maxEvents int
}

// NewInMemoryEventBus creates a new in-memory event bus
func NewInMemoryEventBus(maxEvents int) *InMemoryEventBus {
	return &InMemoryEventBus{
		handlers:  make(map[EventType][]EventHandler),
		events:    make([]*Event, 0, maxEvents),
		maxEvents: maxEvents,
	}
}

// Publish publishes an event
func (b *InMemoryEventBus) Publish(ctx context.Context, event *Event) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	// Store event (with size limit)
	if len(b.events) >= b.maxEvents {
		b.events = b.events[1:]
	}
	b.events = append(b.events, event)

	// Notify handlers
	handlers := b.handlers[event.Type]
	b.mu.Unlock()

	// Call handlers asynchronously
	for _, handler := range handlers {
		go func(h EventHandler) {
			if err := h(ctx, event); err != nil {
				// Log error (in production, use proper logging)
				_ = err
			}
		}(handler)
	}

	b.mu.Lock()
	return nil
}

// Subscribe subscribes to events of a specific type
func (b *InMemoryEventBus) Subscribe(ctx context.Context, eventType EventType, handler EventHandler) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.handlers[eventType] == nil {
		b.handlers[eventType] = make([]EventHandler, 0)
	}

	b.handlers[eventType] = append(b.handlers[eventType], handler)
	return nil
}

// Unsubscribe unsubscribes from events
func (b *InMemoryEventBus) Unsubscribe(eventType EventType, handler EventHandler) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	handlers := b.handlers[eventType]
	for i, h := range handlers {
		// Compare function pointers (simplified - in production use better comparison)
		if &h == &handler {
			b.handlers[eventType] = append(handlers[:i], handlers[i+1:]...)
			return nil
		}
	}

	return nil
}

// GetEvents returns recent events
func (b *InMemoryEventBus) GetEvents(eventType EventType, limit int) []*Event {
	b.mu.RLock()
	defer b.mu.RUnlock()

	var filtered []*Event
	count := 0

	for i := len(b.events) - 1; i >= 0 && count < limit; i-- {
		if eventType == "" || b.events[i].Type == eventType {
			filtered = append(filtered, b.events[i])
			count++
		}
	}

	// Reverse to get chronological order
	for i, j := 0, len(filtered)-1; i < j; i, j = i+1, j-1 {
		filtered[i], filtered[j] = filtered[j], filtered[i]
	}

	return filtered
}

