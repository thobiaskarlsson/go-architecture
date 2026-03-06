package domain

import "time"

type VehicleStatus struct {
	VehicleID  string
	EngineOn   bool
	SpeedKPH   float64
	FuelPct    float64
	CapturedAt time.Time
}

type VehiclePosition struct {
	VehicleID  string
	Latitude   float64
	Longitude  float64
	HeadingDeg float64
	CapturedAt time.Time
}

type WarningSeverity string

const (
	WarningSeverityInfo     WarningSeverity = "info"
	WarningSeverityWarning  WarningSeverity = "warning"
	WarningSeverityCritical WarningSeverity = "critical"
)

type VehicleWarning struct {
	VehicleID  string
	Code       string
	Severity   WarningSeverity
	Message    string
	CapturedAt time.Time
}
