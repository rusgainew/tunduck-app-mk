package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rusgainew/tunduck-app-mk/auth-service/internal/domain/entity"
)

// EventPublisherRabbitMQ - Реализация EventPublisher для RabbitMQ
type EventPublisherRabbitMQ struct {
	ch *amqp.Channel
}

// NewEventPublisherRabbitMQ - Factory
func NewEventPublisherRabbitMQ(ch *amqp.Channel) *EventPublisherRabbitMQ {
	return &EventPublisherRabbitMQ{ch: ch}
}

// PublishUserRegistered - публикует событие регистрации
func (p *EventPublisherRabbitMQ) PublishUserRegistered(ctx context.Context, user *entity.User) error {
	event := map[string]interface{}{
		"event_type": "user.registered",
		"user_id":    user.ID,
		"email":      user.Email,
		"name":       user.Name,
		"timestamp":  user.CreatedAt,
	}

	return p.publishEvent(ctx, "tunduck.auth", "user.registered", event)
}

// PublishUserLoggedIn - публикует событие входа
func (p *EventPublisherRabbitMQ) PublishUserLoggedIn(ctx context.Context, userID string) error {
	event := map[string]interface{}{
		"event_type": "user.logged_in",
		"user_id":    userID,
		"timestamp":  time.Now(),
	}

	return p.publishEvent(ctx, "tunduck.auth", "user.logged_in", event)
}

// PublishUserLoggedOut - публикует событие выхода
func (p *EventPublisherRabbitMQ) PublishUserLoggedOut(ctx context.Context, userID string) error {
	event := map[string]interface{}{
		"event_type": "user.logged_out",
		"user_id":    userID,
		"timestamp":  time.Now(),
	}

	return p.publishEvent(ctx, "tunduck.auth", "user.logged_out", event)
}

// publishEvent - публикует событие на exchange
func (p *EventPublisherRabbitMQ) publishEvent(
	ctx context.Context,
	exchange string,
	routingKey string,
	event interface{},
) error {
	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	return p.ch.PublishWithContext(
		ctx,
		exchange,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}
