package user

import "time"

// DomainEvent interface that all domain events must implement
type DomainEvent interface {
	EventType() string
	OccurredAt() time.Time
}

// UserRegisteredEvent is fired when a new user registers
type UserRegisteredEvent struct {
	UserId    string
	Email     string
	OccurredAt_ time.Time
}

func NewUserRegisteredEvent(userID, email string) *UserRegisteredEvent {
	return &UserRegisteredEvent{
		UserId:      userID,
		Email:       email,
		OccurredAt_: time.Now(),
	}
}

func (e *UserRegisteredEvent) EventType() string {
	return "user.registered"
}

func (e *UserRegisteredEvent) OccurredAt() time.Time {
	return e.OccurredAt_
}

// UserLoggedInEvent is fired when a user successfully logs in
type UserLoggedInEvent struct {
	UserId    string
	OccurredAt_ time.Time
}

func NewUserLoggedInEvent(userID string) *UserLoggedInEvent {
	return &UserLoggedInEvent{
		UserId:      userID,
		OccurredAt_: time.Now(),
	}
}

func (e *UserLoggedInEvent) EventType() string {
	return "user.logged_in"
}

func (e *UserLoggedInEvent) OccurredAt() time.Time {
	return e.OccurredAt_
}

// EventPublisher defines the interface for publishing domain events
type EventPublisher interface {
	Publish(event DomainEvent) error
}
