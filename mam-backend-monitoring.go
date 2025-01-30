// mam-backend.go
package main

import (
    "database/sql"
    "encoding/json"
    "net/http"
    "os"
    "time"

    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
    "github.com/go-chi/cors"
    "github.com/go-redis/redis/v8"
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    "go.uber.org/zap"
    _ "github.com/lib/pq"
)

// Prometheus metrics
var (
    requestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "http_request_duration_seconds",
            Help:    "Duration of HTTP requests.",
            Buckets: prometheus.DefBuckets,
        },
        []string{"handler", "method", "status"},
    )

    activeRequests = prometheus.NewGauge(
        prometheus.GaugeOpts{
            Name: "http_requests_active",
            Help: "Number of active HTTP requests.",
        },
    )

    databaseConnections = prometheus.NewGauge(
        prometheus.GaugeOpts{
            Name: "database_connections_active",
            Help: "Number of active database connections.",
        },
    )

    cacheHits = prometheus.NewCounter(
        prometheus.CounterOpts{
            Name: "cache_hits_total",
            Help: "Total number of cache hits.",
        },
    )

    cacheMisses = prometheus.NewCounter(
        prometheus.CounterOpts{
            Name: "cache_misses_total",
            Help: "Total number of cache misses.",
        },
    )
)

// Service with logger
type Service struct {
    db      *Database
    cache   *redis.Client
    logger  *zap.Logger
    metrics *Metrics
}

// Metrics holds all Prometheus metrics
type Metrics struct {
    requestDuration     *prometheus.HistogramVec
    activeRequests     prometheus.Gauge
    databaseConnections prometheus.Gauge
    cacheHits          prometheus.Counter
    cacheMisses        prometheus.Counter
}

func initMetrics() *Metrics {
    // Register metrics
    prometheus.MustRegister(requestDuration)
    prometheus.MustRegister(activeRequests)
    prometheus.MustRegister(databaseConnections)
    prometheus.MustRegister(cacheHits)
    prometheus.MustRegister(cacheMisses)

    return &Metrics{
        requestDuration:     requestDuration,
        activeRequests:     activeRequests,
        databaseConnections: databaseConnections,
        cacheHits:          cacheHits,
        cacheMisses:        cacheMisses,
    }
}

// Logging middleware
func loggingMiddleware(logger *zap.Logger) func(next http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()
            ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

            // Log request
            logger.Info("request_started",
                zap.String("method", r.Method),
                zap.String("path", r.URL.Path),
                zap.String("remote_addr", r.RemoteAddr),
                zap.String("user_agent", r.UserAgent()),
            )

            // Track active requests
            activeRequests.Inc()
            defer activeRequests.Dec()

            next.ServeHTTP(ww, r)

            // Log response
            duration := time.Since(start)
            logger.Info("request_completed",
                zap.String("method", r.Method),
                zap.String("path", r.URL.Path),
                zap.Int("status", ww.Status()),
                zap.Duration("duration", duration),
                zap.Int("bytes", ww.BytesWritten()),
            )

            // Record metrics
            requestDuration.WithLabelValues(
                r.URL.Path,
                r.Method,
                string(ww.Status()),
            ).Observe(duration.Seconds())
        })
    }
}

// Health check handler with detailed status
func (s *Service) healthHandler(w http.ResponseWriter, r *http.Request) {
    health := struct {
        Status    string `json:"status"`
        Database  bool   `json:"database"`
        Cache     bool   `json:"cache"`
        Timestamp string `json:"timestamp"`
    }{
        Status:    "healthy",
        Timestamp: time.Now().UTC().Format(time.RFC3339),
    }

    // Check database
    if err := s.db.db.Ping(); err != nil {
        health.Status = "degraded"
        health.Database = false
        s.logger.Error("database_health_check_failed", zap.Error(err))
    } else {
        health.Database = true
    }

    // Check cache
    if err := s.cache.Ping(r.Context()).Err(); err != nil {
        health.Status = "degraded"
        health.Cache = false
        s.logger.Error("cache_health_check_failed", zap.Error(err))
    } else {
        health.Cache = true
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(health)
}

func main() {
    // Initialize logger
    logger, err := zap.NewProduction()
    if err != nil {
        panic(err)
    }
    defer logger.Sync()

    // Initialize metrics
    metrics := initMetrics()

    // Initialize other dependencies
    config := loadConfig()
    db := initDatabase(config)
    cache := initRedis(config)

    // Create service instance
    svc := &Service{
        db:      db,
        cache:   cache,
        logger:  logger,
        metrics: metrics,
    }

    // Setup router
    r := chi.NewRouter()

    // Middleware
    r.Use(middleware.Recoverer)
    r.Use(loggingMiddleware(logger))
    r.Use(cors.Handler(cors.Options{
        AllowedOrigins:   []string{"*"},
        AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
        ExposedHeaders:   []string{"Link"},
        AllowCredentials: true,
        MaxAge:           300,
    }))

    // Routes
    r.Route("/api/v1", func(r chi.Router) {
        r.Get("/health", svc.healthHandler)
        // Add more routes here
    })

    // Metrics endpoint
    r.Handle("/metrics", promhttp.Handler())

    // Track database connections
    go func() {
        ticker := time.NewTicker(15 * time.Second)
        for range ticker.C {
            stats := db.db.Stats()
            databaseConnections.Set(float64(stats.OpenConnections))
        }
    }()

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    logger.Info("server_starting", zap.String("port", port))
    if err := http.ListenAndServe(":"+port, r); err != nil {
        logger.Fatal("server_failed",
            zap.Error(err),
        )
    }
}
