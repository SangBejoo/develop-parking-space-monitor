package repository

import (
    "context"
    "github.com/SangBejoo/service-parking-monitor/internal/domain"
)

type ParkingRepository interface {
    GetTrxSupplies(ctx context.Context) ([]domain.TrxSupply, error)
    GetHexagonPlaces(ctx context.Context) ([]domain.MapHexagonPlace, error)
    SaveMonitoring(ctx context.Context, monitoring domain.MonitoringPlace) error
    CreateHexagonPlace(ctx context.Context, hexagon domain.MapHexagonPlace) error
    CreateTrxSupply(ctx context.Context, supply domain.TrxSupply) error
}

type Tile38Repository interface {
    SetLocation(ctx context.Context, fleet string, lat, lon float64) error
    GetLocationsInPolygon(ctx context.Context, polygon []domain.Point) ([]string, error)
}