package nats

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	sharednats "go-architecture/internal/adapters/shared/nats"
	"go-architecture/internal/domain"
	"go-architecture/internal/ports"
)

var ErrUnsupportedSubject = errors.New("unsupported nats subject")

type Subscriber struct {
	app ports.TelemetryInputPort
}

func NewSubscriber(app ports.TelemetryInputPort) *Subscriber {
	return &Subscriber{app: app}
}

type VehicleStatusMessage struct {
	VehicleID string  `json:"vehicle_id"`
	EngineOn  bool    `json:"engine_on"`
	SpeedKPH  float64 `json:"speed_kph"`
	FuelPct   float64 `json:"fuel_pct"`
}

type VehiclePositionMessage struct {
	VehicleID  string  `json:"vehicle_id"`
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
	HeadingDeg float64 `json:"heading_deg"`
}

type VehicleWarningMessage struct {
	VehicleID string                 `json:"vehicle_id"`
	Code      string                 `json:"code"`
	Severity  domain.WarningSeverity `json:"severity"`
	Message   string                 `json:"message"`
}

func (s *Subscriber) Handle(ctx context.Context, subject string, payload []byte) error {
	switch {
	case sharednats.MatchesSubject(subject, sharednats.SubjectVehicleStatus):
		return s.OnVehicleStatus(ctx, payload)
	case sharednats.MatchesSubject(subject, sharednats.SubjectVehiclePosition):
		return s.OnVehiclePosition(ctx, payload)
	case sharednats.MatchesSubject(subject, sharednats.SubjectVehicleWarning):
		return s.OnVehicleWarning(ctx, payload)
	default:
		return ErrUnsupportedSubject
	}
}

func (s *Subscriber) OnVehicleStatus(ctx context.Context, payload []byte) error {
	var msg VehicleStatusMessage
	if err := json.Unmarshal(payload, &msg); err != nil {
		return err
	}

	return s.app.RecordVehicleStatus(ctx, domain.VehicleStatus{
		VehicleID:  msg.VehicleID,
		EngineOn:   msg.EngineOn,
		SpeedKPH:   msg.SpeedKPH,
		FuelPct:    msg.FuelPct,
		CapturedAt: time.Now().UTC(),
	})
}

func (s *Subscriber) OnVehiclePosition(ctx context.Context, payload []byte) error {
	var msg VehiclePositionMessage
	if err := json.Unmarshal(payload, &msg); err != nil {
		return err
	}

	return s.app.RecordVehiclePosition(ctx, domain.VehiclePosition{
		VehicleID:  msg.VehicleID,
		Latitude:   msg.Latitude,
		Longitude:  msg.Longitude,
		HeadingDeg: msg.HeadingDeg,
		CapturedAt: time.Now().UTC(),
	})
}

func (s *Subscriber) OnVehicleWarning(ctx context.Context, payload []byte) error {
	var msg VehicleWarningMessage
	if err := json.Unmarshal(payload, &msg); err != nil {
		return err
	}

	return s.app.RecordVehicleWarning(ctx, domain.VehicleWarning{
		VehicleID:  msg.VehicleID,
		Code:       msg.Code,
		Severity:   msg.Severity,
		Message:    msg.Message,
		CapturedAt: time.Now().UTC(),
	})
}
