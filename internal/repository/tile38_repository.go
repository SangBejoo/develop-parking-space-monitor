package repository

import (
	"context"
	"fmt"
	"github.com/xjem/t38c"
	"github.com/SangBejoo/service-parking-monitor/internal/domain"
)

type tile38Repository struct {
	client *t38c.Client
}

func NewTile38Repository(host string, port int) (*tile38Repository, error) {
	client, err := t38c.New(t38c.Config{
		Address: fmt.Sprintf("%s:%d", host, port),
	})
	if err != nil {
		return nil, err
	}
	return &tile38Repository{client: client}, nil
}

func (r *tile38Repository) SetLocation(ctx context.Context, fleet string, lat, lon float64) error {
    err := r.client.Keys.Set("fleet", fleet).Point(lat, lon).Do(ctx)
    return err
}

func (r *tile38Repository) GetLocationsInPolygon(ctx context.Context, polygon []domain.Point) ([]string, error) {
    points := make([]t38c.Point, len(polygon))
    for i, p := range polygon {
        points[i] = t38c.Point{Lon: p.Longitude, Lat: p.Latitude}
    }

    // Use Within with explicit area search
    search := r.client.Search.Within("fleet")
    resp, err := search.Bounds(
        points[0].Lat, points[0].Lon,
        points[2].Lat, points[2].Lon,
    ).Do(ctx)
    
    if err != nil {
        return []string{}, err // Return empty slice instead of nil
    }
    
    var fleets []string
    for _, obj := range resp.Objects {
        fleets = append(fleets, obj.ID)
    }
    return fleets, nil
}