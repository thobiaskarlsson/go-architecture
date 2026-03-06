package nats

import (
	"context"
	"encoding/json"
	"sync"

	sharednats "go-architecture/internal/adapters/shared/nats"
	"go-architecture/internal/domain"
)

type Message struct {
	Subject string
	Payload []byte
}

type Publisher struct {
	mu       sync.RWMutex
	prefix   string
	messages []Message
}

func NewPublisher(subjectPrefix string) *Publisher {
	return &Publisher{prefix: subjectPrefix}
}

func (p *Publisher) ForwardVehicleStatus(_ context.Context, status domain.VehicleStatus) error {
	return p.publish(sharednats.SubjectVehicleStatus, status)
}

func (p *Publisher) ForwardVehiclePosition(_ context.Context, position domain.VehiclePosition) error {
	return p.publish(sharednats.SubjectVehiclePosition, position)
}

func (p *Publisher) ForwardVehicleWarning(_ context.Context, warning domain.VehicleWarning) error {
	return p.publish(sharednats.SubjectVehicleWarning, warning)
}

func (p *Publisher) publish(subject string, value any) error {
	payload, err := json.Marshal(value)
	if err != nil {
		return err
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	p.messages = append(p.messages, Message{
		Subject: sharednats.ComposeSubject(p.prefix, subject),
		Payload: payload,
	})

	return nil
}

func (p *Publisher) Messages() []Message {
	p.mu.RLock()
	defer p.mu.RUnlock()

	out := make([]Message, len(p.messages))
	for i := range p.messages {
		out[i] = Message{
			Subject: p.messages[i].Subject,
			Payload: append([]byte(nil), p.messages[i].Payload...),
		}
	}
	return out
}
