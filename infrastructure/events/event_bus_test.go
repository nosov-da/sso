package events

import (
	"sso/domain/user"
	"testing"
)

func TestInMemoryEventPublisher_Publish(t *testing.T) {
	publisher := NewInMemoryEventPublisher()

	eventPublished := false
	publisher.Subscribe("user.registered", func(event user.DomainEvent) error {
		eventPublished = true
		return nil
	})

	event := user.NewUserRegisteredEvent("user123", "user@example.com")
	err := publisher.Publish(event)

	if err != nil {
		t.Errorf("Publish() error = %v, want nil", err)
	}

	if !eventPublished {
		t.Error("Event handler should have been called")
	}
}

func TestInMemoryEventPublisher_MultipleSubscribers(t *testing.T) {
	publisher := NewInMemoryEventPublisher()

	calls := 0
	publisher.Subscribe("user.logged_in", func(event user.DomainEvent) error {
		calls++
		return nil
	})

	publisher.Subscribe("user.logged_in", func(event user.DomainEvent) error {
		calls++
		return nil
	})

	event := user.NewUserLoggedInEvent("user123")
	publisher.Publish(event)

	if calls != 2 {
		t.Errorf("Expected 2 handler calls, got %d", calls)
	}
}

func TestInMemoryEventPublisher_NoSubscribers(t *testing.T) {
	publisher := NewInMemoryEventPublisher()

	event := user.NewUserRegisteredEvent("user123", "user@example.com")
	err := publisher.Publish(event)

	if err != nil {
		t.Errorf("Publish() with no subscribers should not error: %v", err)
	}
}

func TestInMemoryEventPublisher_DifferentEventTypes(t *testing.T) {
	publisher := NewInMemoryEventPublisher()

	registeredCalled := false
	loggedInCalled := false

	publisher.Subscribe("user.registered", func(event user.DomainEvent) error {
		registeredCalled = true
		return nil
	})

	publisher.Subscribe("user.logged_in", func(event user.DomainEvent) error {
		loggedInCalled = true
		return nil
	})

	registeredEvent := user.NewUserRegisteredEvent("user123", "user@example.com")
	publisher.Publish(registeredEvent)

	if !registeredCalled || loggedInCalled {
		t.Error("Only registered handler should be called")
	}

	registeredCalled = false
	loggedInEvent := user.NewUserLoggedInEvent("user123")
	publisher.Publish(loggedInEvent)

	if registeredCalled || !loggedInCalled {
		t.Error("Only logged_in handler should be called")
	}
}

func TestNullEventPublisher_Publish(t *testing.T) {
	publisher := &NullEventPublisher{}

	event := user.NewUserRegisteredEvent("user123", "user@example.com")
	err := publisher.Publish(event)

	if err != nil {
		t.Errorf("NullEventPublisher.Publish() error = %v, want nil", err)
	}
	// No events should be published, but it should not error
}
