package events

import (
	"sso/domain/user"
)

// InMemoryEventPublisher publishes events in-memory
type InMemoryEventPublisher struct {
	subscribers map[string][]EventHandler
}

// EventHandler is a function that handles a domain event
type EventHandler func(event user.DomainEvent) error

// NewInMemoryEventPublisher creates a new event publisher
func NewInMemoryEventPublisher() *InMemoryEventPublisher {
	return &InMemoryEventPublisher{
		subscribers: make(map[string][]EventHandler),
	}
}

// Publish publishes an event to all subscribers
func (p *InMemoryEventPublisher) Publish(event user.DomainEvent) error {
	eventType := event.EventType()

	handlers, exists := p.subscribers[eventType]
	if !exists {
		return nil
	}

	for _, handler := range handlers {
		if err := handler(event); err != nil {
			// Log but don't stop other handlers
			_ = err
		}
	}

	return nil
}

// Subscribe registers a handler for an event type
func (p *InMemoryEventPublisher) Subscribe(eventType string, handler EventHandler) {
	if p.subscribers == nil {
		p.subscribers = make(map[string][]EventHandler)
	}

	p.subscribers[eventType] = append(p.subscribers[eventType], handler)
}

// NullEventPublisher is a no-op event publisher (for tests or when events are not needed)
type NullEventPublisher struct{}

// Publish does nothing
func (p *NullEventPublisher) Publish(event user.DomainEvent) error {
	return nil
}
