package http_test

import (
	"context"
	"testing"

	"go-architecture/internal/adapters/inbound/http"
	httpout "go-architecture/internal/adapters/outbound/http"
	sharedhttp "go-architecture/internal/adapters/shared/http"
	"go-architecture/internal/application/telemetry"
	"go-architecture/internal/domain"
)

func TestController_ForwardsTelemetry(t *testing.T) {
	t.Parallel()

	publisher := httpout.NewPublisher("https://northbound.example/telemetry")
	svc := telemetry.NewService(publisher)
	controller := http.NewController(svc)

	if err := controller.PostVehicleStatus(context.Background(), http.VehicleStatusRequest{
		VehicleID: "car-7",
		EngineOn:  true,
		SpeedKPH:  192.0,
		FuelPct:   45.0,
	}); err != nil {
		t.Fatalf("post status: %v", err)
	}

	if err := controller.PostVehicleWarning(context.Background(), http.VehicleWarningRequest{
		VehicleID: "car-7",
		Code:      "ENGINE_TEMP",
		Severity:  domain.WarningSeverityCritical,
		Message:   "engine temp high",
	}); err != nil {
		t.Fatalf("post warning: %v", err)
	}

	statuses := publisher.SentStatuses()
	if len(statuses) != 1 {
		t.Fatalf("statuses len = %d, want 1", len(statuses))
	}
	if statuses[0].SpeedKPH != 192 {
		t.Fatalf("status speed = %v, want 192", statuses[0].SpeedKPH)
	}

	warnings := publisher.SentWarnings()
	if len(warnings) != 1 {
		t.Fatalf("warnings len = %d, want 1", len(warnings))
	}
	if warnings[0].Severity != domain.WarningSeverityCritical {
		t.Fatalf("warning severity = %q, want %q", warnings[0].Severity, domain.WarningSeverityCritical)
	}

	handled := controller.HandledRequests()
	if len(handled) != 2 {
		t.Fatalf("handled len = %d, want 2", len(handled))
	}
	if handled[0].Path != sharedhttp.PathVehicleStatus {
		t.Fatalf("handled path = %q, want %q", handled[0].Path, sharedhttp.PathVehicleStatus)
	}
	if handled[0].Headers[sharedhttp.HeaderTelemetrySource] != sharedhttp.TelemetrySourceProxy {
		t.Fatalf("handled source header = %q, want %q", handled[0].Headers[sharedhttp.HeaderTelemetrySource], sharedhttp.TelemetrySourceProxy)
	}

	sent := publisher.SentRequests()
	if len(sent) != 2 {
		t.Fatalf("sent requests len = %d, want 2", len(sent))
	}
	if sent[1].Path != sharedhttp.PathVehicleWarning {
		t.Fatalf("sent path = %q, want %q", sent[1].Path, sharedhttp.PathVehicleWarning)
	}
	if sent[1].Headers[sharedhttp.HeaderContentType] != sharedhttp.ContentTypeJSON {
		t.Fatalf("content type = %q, want %q", sent[1].Headers[sharedhttp.HeaderContentType], sharedhttp.ContentTypeJSON)
	}
}
