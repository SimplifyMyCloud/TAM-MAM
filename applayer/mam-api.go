// api/server.go
package api

import (
    "encoding/json"
    "net/http"
    "time"
    
    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
    "github.com/go-chi/cors"
)

type Server struct {
    router  *chi.Mux
    service *core.Service
}

func NewServer(service *core.Service) *Server {
    s := &Server{
        router:  chi.NewRouter(),
        service: service,
    }
    
    s.setupMiddleware()
    s.setupRoutes()
    
    return s
}

func (s *Server) setupMiddleware() {
    s.router.Use(middleware.Logger)
    s.router.Use(middleware.Recoverer)
    s.router.Use(middleware.Timeout(60 * time.Second))
    s.router.Use(cors.Handler(cors.Options{
        AllowedOrigins:   []string{"*"},
        AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
        ExposedHeaders:   []string{"Link"},
        AllowCredentials: true,
        MaxAge:           300,
    }))
}

func (s *Server) setupRoutes() {
    s.router.Route("/api/v1", func(r chi.Router) {
        // Public routes
        r.Post("/login", s.handleLogin)
        
        // Protected routes
        r.Group(func(r chi.Router) {
            r.Use(s.authMiddleware)
            
            // Assets
            r.Route("/assets", func(r chi.Router) {
                r.Get("/", s.handleListAssets)
                r.Post("/", s.handleCreateAsset)
                r.Get("/search", s.handleSearchAssets)
                r.Route("/{id}", func(r chi.Router) {
                    r.Get("/", s.handleGetAsset)
                    r.Put("/", s.handleUpdateAsset)
                    r.Delete("/", s.handleDeleteAsset)
                    r.Get("/content", s.handleGetAssetContent)
                })
            })

            // Collections
            r.Route("/collections", func(r chi.Router) {
                r.Get("/", s.handleListCollections)
                r.Post("/", s.handleCreateCollection)
                r.Route("/{id}", func(r chi.Router) {
                    r.Get("/", s.handleGetCollection)
                    r.Put("/", s.handleUpdateCollection)
                    r.Delete("/", s.handleDeleteCollection)
                })
            })
        })
    })
}

// HTTP response wrapper
type Response struct {
    Success bool        `json:"success"`
    Data    interface{} `json:"data,omitempty"`
    Error   string      `json:"error,omitempty"`
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(data)
}

// Asset handlers
func (s *Server) handleListAssets(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    
    limit := 20 // default limit
    offset := 0
    
    assets, err := s.service.ListAssets(ctx, limit, offset)
    if err != nil {
        writeJSON(w, http.StatusInternalServerError, Response{
            Success: false,
            Error:   "Failed to list assets",
        })
        return
    }
    
    writeJSON(w, http.StatusOK, Response{
        Success: true,
        Data:    assets,
    })
}

func (s *Server) handleSearchAssets(w http.ResponseWriter, r *http.Request) {
    var searchParams core.SearchParams
    if err := json.NewDecoder(r.Body).Decode(&searchParams); err != nil {
        writeJSON(w, http.StatusBadRequest, Response{
            Success: false,
            Error:   "Invalid search parameters",
        })
        return
    }
    
    results, err := s.service.SearchAssets(r.Context(), searchParams)
    if err != nil {
        writeJSON(w, http.StatusInternalServerError, Response{
            Success: false,
            Error:   "Search failed",
        })
        return
    }
    
    writeJSON(w, http.StatusOK, Response{
        Success: true,
        Data:    results,
    })
}

func (s *Server) handleCreateAsset(w http.ResponseWriter, r *http.Request) {
    var req core.IngestRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeJSON(w, http.StatusBadRequest, Response{
            Success: false,
            Error:   "Invalid request body",
        })
        return
    }
    
    asset, err := s.service.IngestAsset(r.Context(), req)
    if err != nil {
        writeJSON(w, http.StatusInternalServerError, Response{
            Success: false,
            Error:   "Failed to create asset",
        })
        return
    }
    
    writeJSON(w, http.StatusCreated, Response{
        Success: true,
        Data:    asset,
    })
}

// Middleware
func (s *Server) authMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        if token == "" {
            writeJSON(w, http.StatusUnauthorized, Response{
                Success: false,
                Error:   "No authorization token provided",
            })
            return
        }
        
        // Validate token and get user
        user, err := s.service.ValidateToken(r.Context(), token)
        if err != nil {
            writeJSON(w, http.StatusUnauthorized, Response{
                Success: false,
                Error:   "Invalid token",
            })
            return
        }
        
        // Add user to context
        ctx := context.WithValue(r.Context(), "user", user)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// Error types
type ErrorResponse struct {
    Code    string `json:"code"`
    Message string `json:"message"`
    Details string `json:"details,omitempty"`
}

var (
    ErrInvalidRequest = &ErrorResponse{
        Code:    "INVALID_REQUEST",
        Message: "The request was invalid",
    }
    
    ErrNotFound = &ErrorResponse{
        Code:    "NOT_FOUND",
        Message: "The requested resource was not found",
    }
    
    ErrUnauthorized = &ErrorResponse{
        Code:    "UNAUTHORIZED",
        Message: "Unauthorized access",
    }
)
