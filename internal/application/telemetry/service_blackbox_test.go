package telemetry_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	httpout "go-architecture/internal/adapters/outbound/http"
	natsout "go-architecture/internal/adapters/outbound/nats"
	sharedhttp "go-architecture/internal/adapters/shared/http"
	sharednats "go-architecture/internal/adapters/shared/nats"
	"go-architecture/internal/application/telemetry"
	"go-architecture/internal/domain"
	"go-architecture/internal/ports"
)

type forwarderFactory struct {
	name            string
	new             func() ports.TelemetryForwarderPort
	assertForwarded func(t *testing.T, forwarder ports.TelemetryForwarderPort)
}

func TestService_ForwardsTelemetryAcrossOutboundAdapters(t *testing.T) {
	t.Parallel()

	forwarders := []forwarderFactory{
		{
			name: "http outbound",
			new: func() ports.TelemetryForwarderPort {
				return httpout.NewPublisher("https://northbound.example/telemetry")
			},
			assertForwarded: func(t *testing.T, forwarder ports.TelemetryForwarderPort) {
				t.Helper()
				httpPublisher, ok := forwarder.(*httpout.Publisher)
				if !ok {
					t.Fatalf("forwarder type = %T, want *http.Publisher", forwarder)
				}

				if len(httpPublisher.SentStatuses()) != 1 {
					t.Fatalf("sent statuses len = %d, want 1", len(httpPublisher.SentStatuses()))
				}
				if len(httpPublisher.SentPositions()) != 1 {
					t.Fatalf("sent positions len = %d, want 1", len(httpPublisher.SentPositions()))
				}
				if len(httpPublisher.SentWarnings()) != 1 {
					t.Fatalf("sent warnings len = %d, want 1", len(httpPublisher.SentWarnings()))
				}
				if httpPublisher.SentRequests()[0].Path != sharedhttp.PathVehicleStatus {
					t.Fatalf("path = %q, want %q", httpPublisher.SentRequests()[0].Path, sharedhttp.PathVehicleStatus)
				}
			},
		},
		{
			name: "nats outbound",
			new: func() ports.TelemetryForwarderPort {
				return natsout.NewPublisher("proxy")
			},
			assertForwarded: func(t *testing.T, forwarder ports.TelemetryForwarderPort) {
				t.Helper()
				natsPublisher, ok := forwarder.(*natsout.Publisher)
				if !ok {
					t.Fatalf("forwarder type = %T, want *nats.Publisher", forwarder)
				}

				messages := natsPublisher.Messages()
				if len(messages) != 3 {
					t.Fatalf("messages len = %d, want 3", len(messages))
				}
				if messages[0].Subject != "proxy."+sharednats.SubjectVehicleStatus {
					t.Fatalf("subject = %q, want %q", messages[0].Subject, "proxy."+sharednats.SubjectVehicleStatus)
				}

				var status domain.VehicleStatus
				if err := json.Unmarshal(messages[0].Payload, &status); err != nil {
					t.Fatalf("unmarshal status: %v", err)
				}
				if status.VehicleID != "car-44" {
					t.Fatalf("vehicleID = %q, want %q", status.VehicleID, "car-44")
				}
			},
		},
	}

	for _, factory := range forwarders {
		factory := factory
		t.Run(factory.name, func(t *testing.T) {
			t.Parallel()

			forwarder := factory.new()
			svc := telemetry.NewService(forwarder)
			ctx := context.Background()

			if err := svc.RecordVehicleStatus(ctx, domain.VehicleStatus{
				VehicleID: "car-44",
				EngineOn:  true,
				SpeedKPH:  209.3,
				FuelPct:   31.4,
			}); err != nil {
				t.Fatalf("record status: %v", err)
			}

			if err := svc.RecordVehiclePosition(ctx, domain.VehiclePosition{
				VehicleID:  "car-44",
				Latitude:   57.7089,
				Longitude:  11.9746,
				HeadingDeg: 87.0,
			}); err != nil {
				t.Fatalf("record position: %v", err)
			}

			if err := svc.RecordVehicleWarning(ctx, domain.VehicleWarning{
				VehicleID: "car-44",
				Code:      "LOW_FUEL",
				Severity:  domain.WarningSeverityWarning,
				Message:   "fuel below 35 percent",
			}); err != nil {
				t.Fatalf("record warning: %v", err)
			}

			factory.assertForwarded(t, forwarder)
		})
	}
}

func TestService_ReturnsValidationErrors(t *testing.T) {
	t.Parallel()

	svc := telemetry.NewService(httpout.NewPublisher("https://northbound.example/telemetry"))
	ctx := context.Background()

	err := svc.RecordVehicleStatus(ctx, domain.VehicleStatus{})
	if !errors.Is(err, telemetry.ErrVehicleIDRequired) {
		t.Fatalf("err = %v, want %v", err, telemetry.ErrVehicleIDRequired)
	}

	err = svc.RecordVehicleWarning(ctx, domain.VehicleWarning{
		VehicleID: "car-12",
	})
	if !errors.Is(err, telemetry.ErrWarningCodeNeeded) {
		t.Fatalf("err = %v, want %v", err, telemetry.ErrWarningCodeNeeded)
	}
}
