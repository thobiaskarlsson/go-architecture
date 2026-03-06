package ports

import (
	"context"

	"go-architecture/internal/domain"
)

type TelemetryForwarderPort interface {
	ForwardVehicleStatus(ctx context.Context, status domain.VehicleStatus) error
	ForwardVehiclePosition(ctx context.Context, position domain.VehiclePosition) error
	ForwardVehicleWarning(ctx context.Context, warning domain.VehicleWarning) error
}
