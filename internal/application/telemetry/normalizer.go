package telemetry

import (
	"errors"
	"math"
	"strings"

	"go-architecture/internal/domain"
)

var (
	ErrLatitudeOutOfRange  = errors.New("latitude is out of range")
	ErrLongitudeOutOfRange = errors.New("longitude is out of range")
)

type normalizer interface {
	NormalizeStatus(status domain.VehicleStatus) (domain.VehicleStatus, error)
	NormalizePosition(position domain.VehiclePosition) (domain.VehiclePosition, error)
	NormalizeWarning(warning domain.VehicleWarning) (domain.VehicleWarning, error)
}

type defaultNormalizer struct{}

func (defaultNormalizer) NormalizeStatus(status domain.VehicleStatus) (domain.VehicleStatus, error) {
	status.VehicleID = strings.TrimSpace(status.VehicleID)
	if status.SpeedKPH < 0 {
		status.SpeedKPH = 0
	}
	if status.FuelPct < 0 {
		status.FuelPct = 0
	}
	if status.FuelPct > 100 {
		status.FuelPct = 100
	}

	return status, nil
}

func (defaultNormalizer) NormalizePosition(position domain.VehiclePosition) (domain.VehiclePosition, error) {
	position.VehicleID = strings.TrimSpace(position.VehicleID)
	if position.Latitude < -90 || position.Latitude > 90 {
		return domain.VehiclePosition{}, ErrLatitudeOutOfRange
	}
	if position.Longitude < -180 || position.Longitude > 180 {
		return domain.VehiclePosition{}, ErrLongitudeOutOfRange
	}

	position.HeadingDeg = math.Mod(position.HeadingDeg, 360)
	if position.HeadingDeg < 0 {
		position.HeadingDeg += 360
	}

	return position, nil
}

func (defaultNormalizer) NormalizeWarning(warning domain.VehicleWarning) (domain.VehicleWarning, error) {
	warning.VehicleID = strings.TrimSpace(warning.VehicleID)
	warning.Code = strings.ToUpper(strings.TrimSpace(warning.Code))
	warning.Message = strings.TrimSpace(warning.Message)

	return warning, nil
}
