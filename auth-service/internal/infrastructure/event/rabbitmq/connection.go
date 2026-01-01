package rabbitmq

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Config - конфигурация для RabbitMQ
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	VHost    string
}

// NewConnection - создает новое подключение к RabbitMQ
func NewConnection(cfg Config) (*amqp.Connection, error) {
	url := fmt.Sprintf(
		"amqp://%s:%s@%s:%d/%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.VHost,
	)

	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	return conn, nil
}

// NewChannel - создает новый канал
func NewChannel(conn *amqp.Connection) (*amqp.Channel, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	return ch, nil
}

// DeclareExchange - объявляет exchange
func DeclareExchange(ch *amqp.Channel, name string) error {
	err := ch.ExchangeDeclare(
		name,    // name
		"topic", // type
		true,    // durable
		false,   // auto-deleted
		false,   // internal
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare exchange: %w", err)
	}

	return nil
}

// DefaultConfig - конфигурация по умолчанию
func DefaultConfig() Config {
	return Config{
		Host:     "localhost",
		Port:     5672,
		User:     "guest",
		Password: "guest",
		VHost:    "/",
	}
}
