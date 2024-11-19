package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "os/signal"
    "syscall"
    "database/sql"
    _ "github.com/lib/pq"
    "net/http"
    "github.com/SangBejoo/service-parking-monitor/internal/config"
    "github.com/SangBejoo/service-parking-monitor/internal/repository"
    "github.com/SangBejoo/service-parking-monitor/internal/usecase"
    "github.com/SangBejoo/service-parking-monitor/internal/delivery/http/middleware"
    delivery "github.com/SangBejoo/service-parking-monitor/internal/delivery/http"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    "time"
    "sync"
)

func main() {
    cfg, err := config.LoadConfig()
    if err != nil {
        log.Fatal(err)
    }

    // Setup database with retry
    var db *sql.DB
    maxRetries := 5
    for i := 0; i < maxRetries; i++ {
        connStr := fmt.Sprintf(
            "host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
            cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
        )
        log.Printf("Attempt %d/%d: Connecting to database...", i+1, maxRetries)
        
        db, err = sql.Open("postgres", connStr)
        if err != nil {
            log.Printf("Failed to open database connection: %v", err)
            time.Sleep(5 * time.Second)
            continue
        }
        
        err = db.Ping()
        if err != nil {
            log.Printf("Failed to ping database: %v", err)
            db.Close()
            time.Sleep(5 * time.Second)
            continue
        }
        
        log.Println("Successfully connected to database")
        break
    }
    if err != nil {
        log.Fatalf("Could not connect to the database after retries: %v", err)
    }

    defer db.Close()

    // Initialize repositories
    sqlRepo := repository.NewSQLRepository(db)
    tile38Repo, err := repository.NewTile38Repository(cfg.Tile38Host, cfg.Tile38Port)
    if err != nil {
        log.Fatal(err)
    }
    _ = repository.NewRedisRepository(fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort))

    // Setup monitoring use case and HTTP handler
    monitoringUseCase := usecase.NewMonitoringUseCase(sqlRepo, tile38Repo)
    handler := delivery.NewHandler(monitoringUseCase)

    // Setup HTTP router with middleware
    router := http.NewServeMux()
    router.Handle("/monitoring", middleware.LoggingMiddleware(http.HandlerFunc(handler.GetMonitoring)))
    router.Handle("/trx-supply", middleware.LoggingMiddleware(http.HandlerFunc(handler.CreateTrxSupply)))
    router.Handle("/hexagon-place", middleware.LoggingMiddleware(http.HandlerFunc(handler.CreateHexagonPlace)))
    router.Handle("/set-location", middleware.LoggingMiddleware(http.HandlerFunc(handler.SetLocation)))
    router.Handle("/get-locations-in-polygon", middleware.LoggingMiddleware(http.HandlerFunc(handler.GetLocationsInPolygon)))
    router.Handle("/health", middleware.LoggingMiddleware(http.HandlerFunc(handler.HealthCheck)))
    router.Handle("/metrics", promhttp.Handler())

    // Start monitoring goroutine
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    go monitoringUseCase.StartMonitoring(ctx)

    // Create a WaitGroup to manage goroutines
    var wg sync.WaitGroup
    
    // Setup HTTP server
    server := &http.Server{
        Addr:    ":8080",
        Handler: router,
    }

    // Add server goroutine to wait group
    wg.Add(1)
    go func() {
        defer wg.Done()
        log.Printf("Starting server on %s", server.Addr)
        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Printf("Server error: %v", err)
        }
    }()

    // Handle graceful shutdown
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

    // Wait for interrupt signal
    sig := <-sigChan
    log.Printf("Received signal: %v, initiating shutdown", sig)

    // Create a timeout context for shutdown
    shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer shutdownCancel()

    // Shutdown the server gracefully
    if err := server.Shutdown(shutdownCtx); err != nil {
        log.Printf("Server shutdown error: %v", err)
    }

    // Cancel the monitoring context
    cancel()

    // Wait for all goroutines to finish
    wg.Wait()
    log.Println("Server shutdown completed")
}