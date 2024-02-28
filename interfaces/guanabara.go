package interfaces

import (
	"autflow_back/models"
	"context"
)

type GuanabaraClientRepository interface {
	GetLocations(ctx context.Context) ([]models.Location, []models.Group, error)
	GetTrips(ctx context.Context, d models.GetTrips) ([]models.TripInfoResponse, error)
	GetCity(ctx context.Context, arguments string) ([]string, error)
	GetSchedules(ctx context.Context, departure, arrival, dataDeparture string) ([]models.TripInfoResponse, error)
	RouteValidation(ctx context.Context, departure, arrival string) (bool, error)
}
