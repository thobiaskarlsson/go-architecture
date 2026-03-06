package http

import (
	"context"
	"sync"
	"time"

	sharedhttp "go-architecture/internal/adapters/shared/http"
	"go-architecture/internal/domain"
	"go-architecture/internal/ports"
)

type HandledRequest struct {
	Path    string
	Headers map[string]string
}

type Controller struct {
	app     ports.TelemetryInputPort
	mu      sync.RWMutex
	handled []HandledRequest
}

func NewController(app ports.TelemetryInputPort) *Controller {
	return &Controller{app: app}
}

type VehicleStatusRequest struct {
	VehicleID string
	EngineOn  bool
	SpeedKPH  float64
	FuelPct   float64
}

type VehiclePositionRequest struct {
	VehicleID  string
	Latitude   float64
	Longitude  float64
	HeadingDeg float64
}

type VehicleWarningRequest struct {
	VehicleID string
	Code      string
	Severity  domain.WarningSeverity
	Message   string
}

func (c *Controller) PostVehicleStatus(ctx context.Context, req VehicleStatusRequest) error {
	c.trackRequest(sharedhttp.PathVehicleStatus, req.VehicleID)

	return c.app.RecordVehicleStatus(ctx, domain.VehicleStatus{
		VehicleID:  req.VehicleID,
		EngineOn:   req.EngineOn,
		SpeedKPH:   req.SpeedKPH,
		FuelPct:    req.FuelPct,
		CapturedAt: time.Now().UTC(),
	})
}

func (c *Controller) PostVehiclePosition(ctx context.Context, req VehiclePositionRequest) error {
	c.trackRequest(sharedhttp.PathVehiclePosition, req.VehicleID)

	return c.app.RecordVehiclePosition(ctx, domain.VehiclePosition{
		VehicleID:  req.VehicleID,
		Latitude:   req.Latitude,
		Longitude:  req.Longitude,
		HeadingDeg: req.HeadingDeg,
		CapturedAt: time.Now().UTC(),
	})
}

func (c *Controller) PostVehicleWarning(ctx context.Context, req VehicleWarningRequest) error {
	c.trackRequest(sharedhttp.PathVehicleWarning, req.VehicleID)

	return c.app.RecordVehicleWarning(ctx, domain.VehicleWarning{
		VehicleID:  req.VehicleID,
		Code:       req.Code,
		Severity:   req.Severity,
		Message:    req.Message,
		CapturedAt: time.Now().UTC(),
	})
}

func (c *Controller) HandledRequests() []HandledRequest {
	c.mu.RLock()
	defer c.mu.RUnlock()

	out := make([]HandledRequest, len(c.handled))
	for i := range c.handled {
		headersCopy := make(map[string]string, len(c.handled[i].Headers))
		for key, value := range c.handled[i].Headers {
			headersCopy[key] = value
		}
		out[i] = HandledRequest{
			Path:    c.handled[i].Path,
			Headers: headersCopy,
		}
	}

	return out
}

func (c *Controller) trackRequest(path, vehicleID string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.handled = append(c.handled, HandledRequest{
		Path: path,
		Headers: map[string]string{
			sharedhttp.HeaderContentType:     sharedhttp.ContentTypeJSON,
			sharedhttp.HeaderVehicleID:       vehicleID,
			sharedhttp.HeaderTelemetrySource: sharedhttp.TelemetrySourceProxy,
		},
	})
}
