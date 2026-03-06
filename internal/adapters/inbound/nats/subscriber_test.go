package nats_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"go-architecture/internal/adapters/inbound/nats"
	natsout "go-architecture/internal/adapters/outbound/nats"
	sharednats "go-architecture/internal/adapters/shared/nats"
	"go-architecture/internal/application/telemetry"
	"go-architecture/internal/domain"
)

func TestSubscriber_ParsesAndForwardsTelemetry(t *testing.T) {
	t.Parallel()

	publisher := natsout.NewPublisher("proxy")
	svc := telemetry.NewService(publisher)
	sub := nats.NewSubscriber(svc)

	statusPayload := []byte(`{"vehicle_id":"car-99","engine_on":true,"speed_kph":120.5,"fuel_pct":88}`)
	if err := sub.Handle(context.Background(), sharednats.SubjectVehicleStatus, statusPayload); err != nil {
		t.Fatalf("status: %v", err)
	}

	warningPayload := []byte(`{"vehicle_id":"car-99","code":"TIRE_TEMP","severity":"warning","message":"rear left high"}`)
	if err := sub.Handle(context.Background(), "proxy."+sharednats.SubjectVehicleWarning, warningPayload); err != nil {
		t.Fatalf("warning: %v", err)
	}

	messages := publisher.Messages()
	if len(messages) != 2 {
		t.Fatalf("messages len = %d, want 2", len(messages))
	}
	if messages[0].Subject != "proxy."+sharednats.SubjectVehicleStatus {
		t.Fatalf("subject = %q, want %q", messages[0].Subject, "proxy."+sharednats.SubjectVehicleStatus)
	}

	var status domain.VehicleStatus
	if err := json.Unmarshal(messages[0].Payload, &status); err != nil {
		t.Fatalf("unmarshal status: %v", err)
	}
	if status.SpeedKPH != 120.5 {
		t.Fatalf("speed = %v, want 120.5", status.SpeedKPH)
	}

	err := sub.Handle(context.Background(), "unknown.subject", statusPayload)
	if !errors.Is(err, nats.ErrUnsupportedSubject) {
		t.Fatalf("err = %v, want %v", err, nats.ErrUnsupportedSubject)
	}
}
