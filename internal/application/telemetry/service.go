package telemetry

import (
	"context"
	"errors"
	"strings"
	"time"

	"go-architecture/internal/domain"
	"go-architecture/internal/ports"
)

var (
	ErrVehicleIDRequired = errors.New("vehicleID is required")
	ErrWarningCodeNeeded = errors.New("warning code is required")
)

type Service struct {
	forwarder  ports.TelemetryForwarderPort
	normalizer normalizer
}

type serviceDependencies struct {
	normalizer normalizer
}

func NewService(forwarder ports.TelemetryForwarderPort) *Service {
	return newServiceWithDependencies(forwarder, serviceDependencies{})
}

func newServiceWithDependencies(forwarder ports.TelemetryForwarderPort, deps serviceDependencies) *Service {
	if forwarder == nil {
		panic("forwarder cannot be nil")
	}
	if deps.normalizer == nil {
		deps.normalizer = defaultNormalizer{}
	}
	return &Service{
		forwarder:  forwarder,
		normalizer: deps.normalizer,
	}
}

func (s *Service) RecordVehicleStatus(ctx context.Context, status domain.VehicleStatus) error {
	if status.CapturedAt.IsZero() {
		status.CapturedAt = time.Now().UTC()
	}
	status, err := s.normalizer.NormalizeStatus(status)
	if err != nil {
		return err
	}
	if strings.TrimSpace(status.VehicleID) == "" {
		return ErrVehicleIDRequired
	}

	return s.forwarder.ForwardVehicleStatus(ctx, status)
}

func (s *Service) RecordVehiclePosition(ctx context.Context, position domain.VehiclePosition) error {
	if position.CapturedAt.IsZero() {
		position.CapturedAt = time.Now().UTC()
	}
	position, err := s.normalizer.NormalizePosition(position)
	if err != nil {
		return err
	}
	if strings.TrimSpace(position.VehicleID) == "" {
		return ErrVehicleIDRequired
	}

	return s.forwarder.ForwardVehiclePosition(ctx, position)
}

func (s *Service) RecordVehicleWarning(ctx context.Context, warning domain.VehicleWarning) error {
	if warning.CapturedAt.IsZero() {
		warning.CapturedAt = time.Now().UTC()
	}
	warning, err := s.normalizer.NormalizeWarning(warning)
	if err != nil {
		return err
	}
	if strings.TrimSpace(warning.VehicleID) == "" {
		return ErrVehicleIDRequired
	}
	if strings.TrimSpace(warning.Code) == "" {
		return ErrWarningCodeNeeded
	}

	return s.forwarder.ForwardVehicleWarning(ctx, warning)
}
