package domain

import "time"

type TrxSupply struct {
    ID          int       `json:"id"`
    FleetNumber string    `json:"fleet_number"`
    Latitude    float64   `json:"latitude"`
    Longitude   float64   `json:"longitude"`
    DriverID    string    `json:"driver_id"`
    CreatedAt   time.Time `json:"created_at"`
}

type MapHexagonPlace struct {
    ID           int       `json:"id"`
    HexagonID    string    `json:"hexagon_id"`
    PlaceID      int       `json:"place_id"`
    PlaceTypeID  int       `json:"place_type_id"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
    Polygon      []Point   `json:"polygon"`
}

type MonitoringPlace struct {
    ID       string    `json:"id"`
    Total    int       `json:"total"`
    Polygon  []Point   `json:"polygon"`
    Drivers  []string  `json:"drivers"`
}

type Point struct {
    Latitude  float64 `json:"latitude"`
    Longitude float64 `json:"longitude"`
}