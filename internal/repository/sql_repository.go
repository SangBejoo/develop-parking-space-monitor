package repository

import (
    "context"
    "database/sql"
    "github.com/SangBejoo/service-parking-monitor/internal/domain"
)

type sqlRepository struct {
    db *sql.DB
}

func NewSQLRepository(db *sql.DB) *sqlRepository {
    return &sqlRepository{db: db}
}

func (r *sqlRepository) GetTrxSupplies(ctx context.Context) ([]domain.TrxSupply, error) {
    query := `SELECT id, fleet_number, latitude, longitude, driver_id, created_at FROM trx_supply`
    rows, err := r.db.QueryContext(ctx, query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var supplies []domain.TrxSupply
    for rows.Next() {
        var supply domain.TrxSupply
        if err := rows.Scan(&supply.ID, &supply.FleetNumber, &supply.Latitude, &supply.Longitude, &supply.DriverID, &supply.CreatedAt); err != nil {
            return nil, err
        }
        supplies = append(supplies, supply)
    }
    return supplies, nil
}

func (r *sqlRepository) GetHexagonPlaces(ctx context.Context) ([]domain.MapHexagonPlace, error) {
    query := `SELECT id, hexagon_id, place_id, place_type_id, created_at, updated_at FROM map_hexagon_place`
    rows, err := r.db.QueryContext(ctx, query)
    if (err != nil) {
        return nil, err
    }
    defer rows.Close()

    var hexagons []domain.MapHexagonPlace
    for rows.Next() {
        var hexagon domain.MapHexagonPlace
        if err := rows.Scan(&hexagon.ID, &hexagon.HexagonID, &hexagon.PlaceID, &hexagon.PlaceTypeID, &hexagon.CreatedAt, &hexagon.UpdatedAt); err != nil {
            return nil, err
        }
        hexagons = append(hexagons, hexagon)
    }
    return hexagons, nil
}

func (r *sqlRepository) CreateTrxSupply(ctx context.Context, supply domain.TrxSupply) error {
    query := `
        INSERT INTO trx_supply (fleet_number, latitude, longitude, driver_id)
        VALUES ($1, $2, $3, $4)
    `
    _, err := r.db.ExecContext(ctx, query,
        supply.FleetNumber,
        supply.Latitude,
        supply.Longitude,
        supply.DriverID,
    )
    return err
}

func (r *sqlRepository) CreateHexagonPlace(ctx context.Context, hexagon domain.MapHexagonPlace) error {
    query := `
        INSERT INTO map_hexagon_place (hexagon_id, place_id, place_type_id)
        VALUES ($1, $2, $3)
    `
    _, err := r.db.ExecContext(ctx, query,
        hexagon.HexagonID,
        hexagon.PlaceID,
        hexagon.PlaceTypeID,
    )
    return err
}

func (r *sqlRepository) SaveMonitoring(ctx context.Context, monitoring domain.MonitoringPlace) error {
    // ...existing code...
    return nil
}