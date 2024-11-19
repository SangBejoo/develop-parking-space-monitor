package http

import (
    "net/http"
    "encoding/json"
    "github.com/SangBejoo/service-parking-monitor/internal/domain"
    "github.com/SangBejoo/service-parking-monitor/internal/usecase"
)

// Add Handler struct
type Handler struct {
    monitoringUseCase *usecase.MonitoringUseCase
}

// NewHandler creates a new Handler
func NewHandler(monitoringUseCase *usecase.MonitoringUseCase) *Handler {
    return &Handler{
        monitoringUseCase: monitoringUseCase,
    }
}

func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
    status := map[string]interface{}{
        "status": "UP",
        "database": "UP",
        "redis": "UP", 
        "tile38": "UP",
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "success": true,
        "data": status,
    })
}

func (h *Handler) CreateTrxSupply(w http.ResponseWriter, r *http.Request) {
    var supply domain.TrxSupply
    if err := json.NewDecoder(r.Body).Decode(&supply); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    err := h.monitoringUseCase.CreateTrxSupply(r.Context(), supply)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
}

func (h *Handler) CreateHexagonPlace(w http.ResponseWriter, r *http.Request) {
    var hexagon domain.MapHexagonPlace
    if err := json.NewDecoder(r.Body).Decode(&hexagon); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    err := h.monitoringUseCase.CreateHexagonPlace(r.Context(), hexagon)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
}

func (h *Handler) GetMonitoring(w http.ResponseWriter, r *http.Request) {
    monitoringData, err := h.monitoringUseCase.GetCurrentMonitoring(r.Context())
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "success": true,
        "data": map[string]interface{}{
            "places": monitoringData,
        },
    })
}

func (h *Handler) SetLocation(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Fleet string  `json:"fleet"`
        Lat   float64 `json:"lat"`
        Lon   float64 `json:"lon"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    err := h.monitoringUseCase.SetLocation(r.Context(), req.Fleet, req.Lat, req.Lon)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetLocationsInPolygon(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Polygon []domain.Point `json:"polygon"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    fleets, err := h.monitoringUseCase.GetLocationsInPolygon(r.Context(), req.Polygon)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(map[string]interface{}{
        "fleets": fleets,
    })
}