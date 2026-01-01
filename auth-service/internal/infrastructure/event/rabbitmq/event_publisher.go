package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"

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

// Publish - универсальный метод для публикации domain events
func (p *EventPublisherRabbitMQ) Publish(ctx context.Context, event entity.DomainEvent) error {
	// Создаем структуру сообщения
	message := map[string]interface{}{
		"event_name":   event.EventName(),
		"aggregate_id": event.AggregateID(),
		"occurred_at":  event.OccurredAt(),
		"payload":      event,
	}

	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// Определяем exchange и routing key на основе типа события
	exchange := "tunduck.auth"
	routingKey := event.EventName()

	err = p.ch.PublishWithContext(
		ctx,
		exchange,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
			Headers: amqp.Table{
				"event_name":   event.EventName(),
				"aggregate_id": event.AggregateID(),
			},
		},
	)

	if err != nil {
		return fmt.Errorf("failed to publish event %s: %w", event.EventName(), err)
	}

	return nil
}
