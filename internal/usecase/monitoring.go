package usecase

import (
    "context"
    "time"
    "github.com/SangBejoo/service-parking-monitor/internal/domain"
    "github.com/SangBejoo/service-parking-monitor/internal/repository"
    "github.com/SangBejoo/service-parking-monitor/internal/util"
    "github.com/SangBejoo/service-parking-monitor/internal/metrics"
)

type MonitoringUseCase struct {
    parkingRepo repository.ParkingRepository
    tile38Repo  repository.Tile38Repository
}

func NewMonitoringUseCase(parkingRepo repository.ParkingRepository, tile38Repo repository.Tile38Repository) *MonitoringUseCase {
    return &MonitoringUseCase{
        parkingRepo: parkingRepo,
        tile38Repo: tile38Repo,
    }
}

func (uc *MonitoringUseCase) StartMonitoring(ctx context.Context) {
    ticker := time.NewTicker(5 * time.Minute)
    for {
        select {
        case <-ticker.C:
            uc.updateParkingStatus(ctx)
        case <-ctx.Done():
            return
        }
    }
}

// Add new method to get current monitoring data
func (uc *MonitoringUseCase) GetCurrentMonitoring(ctx context.Context) ([]domain.MonitoringPlace, error) {
    // Get all hexagon places
    hexagons, err := uc.parkingRepo.GetHexagonPlaces(ctx)
    if err != nil {
        return nil, err
    }

    // Get all supplies for reference
    supplies, err := uc.parkingRepo.GetTrxSupplies(ctx)
    if err != nil {
        return nil, err
    }

    var results []domain.MonitoringPlace
    for _, hexagon := range hexagons {
        polygon := util.GetPolygonPoints(hexagon.HexagonID)
        
        // Get vehicles in this polygon
        var driversInPolygon []string
        for _, supply := range supplies {
            // Check if supply point is inside polygon
            if util.IsPointInPolygon(domain.Point{
                Latitude: supply.Latitude,
                Longitude: supply.Longitude,
            }, polygon) {
                driversInPolygon = append(driversInPolygon, supply.FleetNumber)
            }
        }
        
        results = append(results, domain.MonitoringPlace{
            ID:      hexagon.HexagonID,
            Total:   len(driversInPolygon),
            Polygon: polygon,
            Drivers: driversInPolygon,
        })
    }
    
    return results, nil
}

func (uc *MonitoringUseCase) updateParkingStatus(ctx context.Context) error {
    start := time.Now()
    defer func() {
        metrics.UpdateLatency.Observe(time.Since(start).Seconds())
    }()

    supplies, err := uc.parkingRepo.GetTrxSupplies(ctx)
    if (err != nil) {
        return err
    }

    // Update locations in Tile38
    for _, supply := range supplies {
        err = uc.tile38Repo.SetLocation(ctx, supply.FleetNumber, supply.Latitude, supply.Longitude)
        if err != nil {
            return err
        }
    }

    hexagons, err := uc.parkingRepo.GetHexagonPlaces(ctx)
    if err != nil {
        return err
    }

    // Process each hexagon
    for _, hexagon := range hexagons {
        // Get polygon points for hexagon
        polygon := util.GetPolygonPoints(hexagon.HexagonID)
        
        // Get fleet numbers in polygon
        fleets, err := uc.tile38Repo.GetLocationsInPolygon(ctx, polygon)
        if err != nil {
            continue
        }

        metrics.ParkingOccupancy.WithLabelValues(hexagon.HexagonID).Set(float64(len(fleets)))

        monitoring := domain.MonitoringPlace{
            ID:      hexagon.HexagonID,
            Total:   len(fleets),
            Polygon: polygon,
            Drivers: fleets,
        }

        // Save monitoring data
        err = uc.parkingRepo.SaveMonitoring(ctx, monitoring)
        if err != nil {
            continue
        }
    }

    return nil
}

func (uc *MonitoringUseCase) CreateHexagonPlace(ctx context.Context, hexagon domain.MapHexagonPlace) error {
    // Generate polygon points based on center coordinates
    centerLat := -6.200000 // Example center point
    centerLon := 106.816666
    radius := 0.01 // Approximately 1km radius

    // Generate hexagon vertices
    polygon := []domain.Point{
        {Latitude: centerLat + radius, Longitude: centerLon},
        {Latitude: centerLat + (radius * 0.5), Longitude: centerLon + (radius * 0.866)},
        {Latitude: centerLat - (radius * 0.5), Longitude: centerLon + (radius * 0.866)},
        {Latitude: centerLat - radius, Longitude: centerLon},
        {Latitude: centerLat - (radius * 0.5), Longitude: centerLon - (radius * 0.866)},
        {Latitude: centerLat + (radius * 0.5), Longitude: centerLon - (radius * 0.866)},
    }

    hexagon.Polygon = polygon
    return uc.parkingRepo.CreateHexagonPlace(ctx, hexagon)
}

func (uc *MonitoringUseCase) CreateTrxSupply(ctx context.Context, supply domain.TrxSupply) error {
    // Implementation of CreateTrxSupply
    // ...existing code...
    return nil
}

func (uc *MonitoringUseCase) SetLocation(ctx context.Context, fleet string, lat, lon float64) error {
    return uc.tile38Repo.SetLocation(ctx, fleet, lat, lon)
}

func (uc *MonitoringUseCase) GetLocationsInPolygon(ctx context.Context, polygon []domain.Point) ([]string, error) {
    return uc.tile38Repo.GetLocationsInPolygon(ctx, polygon)
}