package http

const (
	PathVehicleStatus   = "/telemetry/status"
	PathVehiclePosition = "/telemetry/position"
	PathVehicleWarning  = "/telemetry/warning"
)

const (
	HeaderContentType     = "Content-Type"
	HeaderVehicleID       = "X-Vehicle-Id"
	HeaderTelemetrySource = "X-Telemetry-Source"
	ContentTypeJSON       = "application/json"
	TelemetrySourceProxy  = "vehicle-telemetry-proxy"
)
