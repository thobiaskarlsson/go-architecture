package main

import (
	"context"
	"fmt"

	httpin "go-architecture/internal/adapters/inbound/http"
	natsin "go-architecture/internal/adapters/inbound/nats"
	httpout "go-architecture/internal/adapters/outbound/http"
	"go-architecture/internal/adapters/outbound/multi"
	natsout "go-architecture/internal/adapters/outbound/nats"
	"go-architecture/internal/application/telemetry"
)

func main() {
	httpPublisher := httpout.NewPublisher("https://northbound.example/telemetry")
	natsPublisher := natsout.NewPublisher("proxy")
	forwarder := multi.NewForwarder(httpPublisher, natsPublisher)
	service := telemetry.NewService(forwarder)

	httpController := httpin.NewController(service)
	natsSubscriber := natsin.NewSubscriber(service)

	ctx := context.Background()

	_ = httpController.PostVehicleStatus(ctx, httpin.VehicleStatusRequest{
		VehicleID: "car-1",
		EngineOn:  true,
		SpeedKPH:  80,
		FuelPct:   92,
	})

	_ = natsSubscriber.OnVehicleWarning(ctx, []byte(`{"vehicle_id":"car-1","code":"LOW_FUEL","severity":"warning","message":"fuel is dropping"}`))

	fmt.Printf("http outbound status events: %d\n", len(httpPublisher.SentStatuses()))
	fmt.Printf("http outbound warning events: %d\n", len(httpPublisher.SentWarnings()))
	fmt.Printf("nats outbound messages: %d\n", len(natsPublisher.Messages()))
}
