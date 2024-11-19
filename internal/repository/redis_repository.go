
package repository

import (
    "context"
    "encoding/json"
    "github.com/go-redis/redis/v8"
    "github.com/SangBejoo/service-parking-monitor/internal/domain"
)

type redisRepository struct {
    client *redis.Client
}

func NewRedisRepository(addr string) *redisRepository {
    client := redis.NewClient(&redis.Options{
        Addr: addr,
    })
    return &redisRepository{client: client}
}

func (r *redisRepository) SaveMonitoring(ctx context.Context, monitoring domain.MonitoringPlace) error {
    data, err := json.Marshal(monitoring)
    if err != nil {
        return err
    }
    return r.client.Set(ctx, "monitoring:"+monitoring.ID, data, 0).Err()
}