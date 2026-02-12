package messaging

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"log"
	"sync"

	"github.com/erlint1212/portfolio/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	channel *amqp.Channel
	mu      sync.Mutex 
}

func NewPublisher(conn *amqp.Connection) (*Publisher, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	return &Publisher{channel: ch}, nil
}

func (p *Publisher) Close() error {
	return p.channel.Close()
}

func PublishGob[T any](ctx context.Context,ch *amqp.Channel, exchange, key string, val T) error {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(val)
	if err != nil {
		return err
	}

	msg := amqp.Publishing{
		ContentType: "application/gob",
		Body:        buffer.Bytes(),
	}

	return ch.PublishWithContext(ctx, exchange, key, false, false, msg)
}

func (p *Publisher) PublishGameLog(ctx context.Context, gl routing.GameLog) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	key := fmt.Sprintf("%s.guest", routing.GameLogSlug)

	err := PublishGob(
		ctx,
		p.channel,
		routing.ExchangePortfolioTopic,
		key,
		gl,
	)

	if err != nil {
		log.Printf("[ERROR] Error during publishing GameLog: %v", err)
		return err
	}
	return nil
}
