package multi

import (
	"context"

	"go-architecture/internal/domain"
	"go-architecture/internal/ports"
)

type Forwarder struct {
	forwarders []ports.TelemetryForwarderPort
}

func NewForwarder(forwarders ...ports.TelemetryForwarderPort) *Forwarder {
	return &Forwarder{forwarders: forwarders}
}

func (f *Forwarder) ForwardVehicleStatus(ctx context.Context, status domain.VehicleStatus) error {
	for i := range f.forwarders {
		if err := f.forwarders[i].ForwardVehicleStatus(ctx, status); err != nil {
			return err
		}
	}

	return nil
}

func (f *Forwarder) ForwardVehiclePosition(ctx context.Context, position domain.VehiclePosition) error {
	for i := range f.forwarders {
		if err := f.forwarders[i].ForwardVehiclePosition(ctx, position); err != nil {
			return err
		}
	}

	return nil
}

func (f *Forwarder) ForwardVehicleWarning(ctx context.Context, warning domain.VehicleWarning) error {
	for i := range f.forwarders {
		if err := f.forwarders[i].ForwardVehicleWarning(ctx, warning); err != nil {
			return err
		}
	}

	return nil
}
