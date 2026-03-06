package telemetry

import (
	"context"
	"errors"
	"testing"

	"go-architecture/internal/domain"
)

var errNormalizationFailed = errors.New("normalization failed")

type normalizerMock struct {
	normalizeStatusFn   func(status domain.VehicleStatus) (domain.VehicleStatus, error)
	normalizePositionFn func(position domain.VehiclePosition) (domain.VehiclePosition, error)
	normalizeWarningFn  func(warning domain.VehicleWarning) (domain.VehicleWarning, error)
}

func (m *normalizerMock) NormalizeStatus(status domain.VehicleStatus) (domain.VehicleStatus, error) {
	if m.normalizeStatusFn != nil {
		return m.normalizeStatusFn(status)
	}

	return status, nil
}

func (m *normalizerMock) NormalizePosition(position domain.VehiclePosition) (domain.VehiclePosition, error) {
	if m.normalizePositionFn != nil {
		return m.normalizePositionFn(position)
	}

	return position, nil
}

func (m *normalizerMock) NormalizeWarning(warning domain.VehicleWarning) (domain.VehicleWarning, error) {
	if m.normalizeWarningFn != nil {
		return m.normalizeWarningFn(warning)
	}

	return warning, nil
}

type forwarderSpy struct {
	statuses  []domain.VehicleStatus
	positions []domain.VehiclePosition
	warnings  []domain.VehicleWarning
}

func (s *forwarderSpy) ForwardVehicleStatus(_ context.Context, status domain.VehicleStatus) error {
	s.statuses = append(s.statuses, status)
	return nil
}

func (s *forwarderSpy) ForwardVehiclePosition(_ context.Context, position domain.VehiclePosition) error {
	s.positions = append(s.positions, position)
	return nil
}

func (s *forwarderSpy) ForwardVehicleWarning(_ context.Context, warning domain.VehicleWarning) error {
	s.warnings = append(s.warnings, warning)
	return nil
}

func TestService_UsesInjectedNormalizer(t *testing.T) {
	t.Parallel()

	mock := &normalizerMock{
		normalizeStatusFn: func(status domain.VehicleStatus) (domain.VehicleStatus, error) {
			status.VehicleID = "car-normalized"
			status.FuelPct = 100
			status.SpeedKPH = 0
			return status, nil
		},
	}
	spy := &forwarderSpy{}

	svc := newServiceWithDependencies(spy, serviceDependencies{
		normalizer: mock,
	})
	ctx := context.Background()

	if err := svc.RecordVehicleStatus(ctx, domain.VehicleStatus{
		VehicleID: " car-raw ",
		FuelPct:   190,
		SpeedKPH:  -10,
	}); err != nil {
		t.Fatalf("record status: %v", err)
	}

	if len(spy.statuses) != 1 {
		t.Fatalf("statuses len = %d, want 1", len(spy.statuses))
	}
	if spy.statuses[0].VehicleID != "car-normalized" {
		t.Fatalf("vehicleID = %q, want %q", spy.statuses[0].VehicleID, "car-normalized")
	}
	if spy.statuses[0].FuelPct != 100 {
		t.Fatalf("fuelPct = %v, want 100", spy.statuses[0].FuelPct)
	}
	if spy.statuses[0].SpeedKPH != 0 {
		t.Fatalf("speed = %v, want 0", spy.statuses[0].SpeedKPH)
	}
}

func TestService_PropagatesNormalizerErrors(t *testing.T) {
	t.Parallel()

	mock := &normalizerMock{
		normalizePositionFn: func(position domain.VehiclePosition) (domain.VehiclePosition, error) {
			return domain.VehiclePosition{}, errNormalizationFailed
		},
	}
	spy := &forwarderSpy{}

	svc := newServiceWithDependencies(spy, serviceDependencies{
		normalizer: mock,
	})
	err := svc.RecordVehiclePosition(context.Background(), domain.VehiclePosition{
		VehicleID: "car-8",
	})
	if !errors.Is(err, errNormalizationFailed) {
		t.Fatalf("err = %v, want %v", err, errNormalizationFailed)
	}
	if len(spy.positions) != 0 {
		t.Fatalf("positions len = %d, want 0", len(spy.positions))
	}
}

func TestNewServiceWithNilNormalizer_UsesDefault(t *testing.T) {
	t.Parallel()

	spy := &forwarderSpy{}
	svc := newServiceWithDependencies(spy, serviceDependencies{})

	err := svc.RecordVehicleStatus(context.Background(), domain.VehicleStatus{
		VehicleID: "car-1",
		SpeedKPH:  -22,
		FuelPct:   140,
	})
	if err != nil {
		t.Fatalf("record status: %v", err)
	}

	if len(spy.statuses) != 1 {
		t.Fatalf("statuses len = %d, want 1", len(spy.statuses))
	}
	if spy.statuses[0].SpeedKPH != 0 {
		t.Fatalf("speed = %v, want 0", spy.statuses[0].SpeedKPH)
	}
	if spy.statuses[0].FuelPct != 100 {
		t.Fatalf("fuelPct = %v, want 100", spy.statuses[0].FuelPct)
	}
}
