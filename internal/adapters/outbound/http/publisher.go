package http

import (
	"context"
	"sync"

	sharedhttp "go-architecture/internal/adapters/shared/http"
	"go-architecture/internal/domain"
)

type Request struct {
	Path    string
	Headers map[string]string
}

type Publisher struct {
	mu        sync.RWMutex
	baseURL   string
	statuses  []domain.VehicleStatus
	positions []domain.VehiclePosition
	warnings  []domain.VehicleWarning
	requests  []Request
}

func NewPublisher(baseURL string) *Publisher {
	return &Publisher{baseURL: baseURL}
}

func (p *Publisher) ForwardVehicleStatus(_ context.Context, status domain.VehicleStatus) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.requests = append(p.requests, p.newRequest(sharedhttp.PathVehicleStatus, status.VehicleID))
	p.statuses = append(p.statuses, status)
	return nil
}

func (p *Publisher) ForwardVehiclePosition(_ context.Context, position domain.VehiclePosition) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.requests = append(p.requests, p.newRequest(sharedhttp.PathVehiclePosition, position.VehicleID))
	p.positions = append(p.positions, position)
	return nil
}

func (p *Publisher) ForwardVehicleWarning(_ context.Context, warning domain.VehicleWarning) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.requests = append(p.requests, p.newRequest(sharedhttp.PathVehicleWarning, warning.VehicleID))
	p.warnings = append(p.warnings, warning)
	return nil
}

func (p *Publisher) BaseURL() string {
	return p.baseURL
}

func (p *Publisher) SentStatuses() []domain.VehicleStatus {
	p.mu.RLock()
	defer p.mu.RUnlock()

	out := make([]domain.VehicleStatus, len(p.statuses))
	copy(out, p.statuses)
	return out
}

func (p *Publisher) SentPositions() []domain.VehiclePosition {
	p.mu.RLock()
	defer p.mu.RUnlock()

	out := make([]domain.VehiclePosition, len(p.positions))
	copy(out, p.positions)
	return out
}

func (p *Publisher) SentWarnings() []domain.VehicleWarning {
	p.mu.RLock()
	defer p.mu.RUnlock()

	out := make([]domain.VehicleWarning, len(p.warnings))
	copy(out, p.warnings)
	return out
}

func (p *Publisher) SentRequests() []Request {
	p.mu.RLock()
	defer p.mu.RUnlock()

	out := make([]Request, len(p.requests))
	for i := range p.requests {
		headersCopy := make(map[string]string, len(p.requests[i].Headers))
		for key, value := range p.requests[i].Headers {
			headersCopy[key] = value
		}
		out[i] = Request{
			Path:    p.requests[i].Path,
			Headers: headersCopy,
		}
	}

	return out
}

func (p *Publisher) newRequest(path, vehicleID string) Request {
	return Request{
		Path: path,
		Headers: map[string]string{
			sharedhttp.HeaderContentType:     sharedhttp.ContentTypeJSON,
			sharedhttp.HeaderVehicleID:       vehicleID,
			sharedhttp.HeaderTelemetrySource: sharedhttp.TelemetrySourceProxy,
		},
	}
}
