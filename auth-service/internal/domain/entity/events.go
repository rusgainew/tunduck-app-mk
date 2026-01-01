package entity

import "time"

// DomainEvent - Interface for all domain events
type DomainEvent interface {
	EventName() string
	OccurredAt() time.Time
	AggregateID() string
}

// BaseEvent - Common fields for all events
type BaseEvent struct {
	EventID     string
	AggregateId string
	Timestamp   time.Time
}

func (e BaseEvent) OccurredAt() time.Time {
	return e.Timestamp
}

func (e BaseEvent) AggregateID() string {
	return e.AggregateId
}

// UserRegistered - Domain event when user registers
type UserRegistered struct {
	BaseEvent
	UserID    string
	Email     string
	Name      string
	CreatedAt time.Time
}

func (e UserRegistered) EventName() string {
	return "user.registered"
}

// NewUserRegistered - Factory
func NewUserRegistered(userID, email, name string) *UserRegistered {
	return &UserRegistered{
		BaseEvent: BaseEvent{
			AggregateId: userID,
			Timestamp:   time.Now(),
		},
		UserID:    userID,
		Email:     email,
		Name:      name,
		CreatedAt: time.Now(),
	}
}

// UserLoggedIn - Domain event when user logs in
type UserLoggedIn struct {
	BaseEvent
	UserID    string
	Email     string
	LoginAt   time.Time
	IPAddress string
}

func (e UserLoggedIn) EventName() string {
	return "user.logged_in"
}

// NewUserLoggedIn - Factory
func NewUserLoggedIn(userID, email, ipAddress string) *UserLoggedIn {
	return &UserLoggedIn{
		BaseEvent: BaseEvent{
			AggregateId: userID,
			Timestamp:   time.Now(),
		},
		UserID:    userID,
		Email:     email,
		LoginAt:   time.Now(),
		IPAddress: ipAddress,
	}
}

// UserLoggedOut - Domain event when user logs out
type UserLoggedOut struct {
	BaseEvent
	UserID   string
	LogoutAt time.Time
}

func (e UserLoggedOut) EventName() string {
	return "user.logged_out"
}

// NewUserLoggedOut - Factory
func NewUserLoggedOut(userID string) *UserLoggedOut {
	return &UserLoggedOut{
		BaseEvent: BaseEvent{
			AggregateId: userID,
			Timestamp:   time.Now(),
		},
		UserID:   userID,
		LogoutAt: time.Now(),
	}
}

// UserBlocked - Domain event when user is blocked
type UserBlocked struct {
	BaseEvent
	UserID    string
	Reason    string
	BlockedAt time.Time
}

func (e UserBlocked) EventName() string {
	return "user.blocked"
}

// NewUserBlocked - Factory
func NewUserBlocked(userID, reason string) *UserBlocked {
	return &UserBlocked{
		BaseEvent: BaseEvent{
			AggregateId: userID,
			Timestamp:   time.Now(),
		},
		UserID:    userID,
		Reason:    reason,
		BlockedAt: time.Now(),
	}
}

// PasswordChanged - Domain event when password is changed
type PasswordChanged struct {
	BaseEvent
	UserID    string
	ChangedAt time.Time
}

func (e PasswordChanged) EventName() string {
	return "user.password_changed"
}

// NewPasswordChanged - Factory
func NewPasswordChanged(userID string) *PasswordChanged {
	return &PasswordChanged{
		BaseEvent: BaseEvent{
			AggregateId: userID,
			Timestamp:   time.Now(),
		},
		UserID:    userID,
		ChangedAt: time.Now(),
	}
}
