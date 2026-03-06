package ports

import (
	"context"

	"go-architecture/internal/domain"
)

type TelemetryInputPort interface {
	RecordVehicleStatus(ctx context.Context, status domain.VehicleStatus) error
	RecordVehiclePosition(ctx context.Context, position domain.VehiclePosition) error
	RecordVehicleWarning(ctx context.Context, warning domain.VehicleWarning) error
}
