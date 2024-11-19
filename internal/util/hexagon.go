package util

import (
    "math"
    "github.com/SangBejoo/service-parking-monitor/internal/domain"
)

const (
    hexagonSize = 0.01 // Roughly 1km at equator
)

func GetPolygonPoints(hexID string) []domain.Point {
    // Calculate center point from hexagon ID
    lat, lon := getHexagonCenter()
    
    // Generate 6 points of hexagon
    points := make([]domain.Point, 6)
    for i := 0; i < 6; i++ {
        angle := float64(i) * (math.Pi / 3)
        points[i] = domain.Point{
            Latitude:  lat + hexagonSize*math.Sin(angle),
            Longitude: lon + hexagonSize*math.Cos(angle),
        }
    }
    return points
}

func getHexagonCenter() (float64, float64) {
    // In real implementation, decode hexID to get center coordinates
    // This is simplified example returning fixed point
    return -6.2, 106.8
}

// IsPointInPolygon checks if a point is inside a polygon using ray casting algorithm
func IsPointInPolygon(point domain.Point, polygon []domain.Point) bool {
    inside := false
    j := len(polygon) - 1
    
    for i := 0; i < len(polygon); i++ {
        if ((polygon[i].Latitude > point.Latitude) != (polygon[j].Latitude > point.Latitude)) &&
            (point.Longitude < (polygon[j].Longitude-polygon[i].Longitude)*(point.Latitude-polygon[i].Latitude)/
                (polygon[j].Latitude-polygon[i].Latitude)+polygon[i].Longitude) {
            inside = !inside
        }
        j = i
    }
    
    return inside
}