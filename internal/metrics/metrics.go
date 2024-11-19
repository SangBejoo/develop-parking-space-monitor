
package metrics

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    UpdateLatency = promauto.NewHistogram(prometheus.HistogramOpts{
        Name: "parking_monitor_update_duration_seconds",
        Help: "Time taken to update parking status",
    })

    ParkingOccupancy = promauto.NewGaugeVec(prometheus.GaugeOpts{
        Name: "parking_monitor_occupancy",
        Help: "Number of vehicles in parking area",
    }, []string{"hexagon_id"})
)