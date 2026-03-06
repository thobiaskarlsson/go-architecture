package nats

import "strings"

const (
	SubjectVehicleStatus   = "telemetry.vehicle.status"
	SubjectVehiclePosition = "telemetry.vehicle.position"
	SubjectVehicleWarning  = "telemetry.vehicle.warning"
)

func ComposeSubject(prefix, subject string) string {
	if prefix == "" {
		return subject
	}

	return prefix + "." + subject
}

func MatchesSubject(subject, expected string) bool {
	return subject == expected || strings.HasSuffix(subject, "."+expected)
}
