// core/service.go
package core

import (
    "context"
    "encoding/json"
    "fmt"
    "time"
    
    "github.com/go-redis/redis/v8"
)

type Service struct {
    db         *Database
    cache      *redis.Client
    tams       *TamsClient
    processor  *MediaProcessor
}

type Config struct {
    TamsBaseURL     string
    CacheExpiration time.Duration
}

// New creates a new instance of the core service
func New(db *Database, cache *redis.Client, config Config) *Service {
    return &Service{
        db:        db,
        cache:     cache,
        tams:      NewTamsClient(config.TamsBaseURL),
        processor: NewMediaProcessor(),
    }
}

// Asset workflow states
const (
    StateNew       = "new"
    StateIngesting = "ingesting"
    StateProcessing = "processing"
    StateReady     = "ready"
    StateFailed    = "failed"
)

// IngestAsset handles the complete asset ingest workflow
func (s *Service) IngestAsset(ctx context.Context, req IngestRequest) (*Asset, error) {
    // Create initial asset record
    asset := &Asset{
        Title:       req.Title,
        Description: req.Description,
        Type:        req.Type,
        Status:      StateNew,
        Metadata:    req.Metadata,
        CreatedBy:   req.UserID,
    }

    // Save to database
    if err := s.db.CreateAsset(ctx, asset); err != nil {
        return nil, fmt.Errorf("failed to create asset: %w", err)
    }

    // Start async ingest process
    go s.processIngest(context.Background(), asset.ID, req)

    return asset, nil
}

type IngestRequest struct {
    Title       string
    Description string
    Type        string
    UserID      string
    Metadata    json.RawMessage
    SourcePath  string
}

// processIngest handles the async ingest workflow
func (s *Service) processIngest(ctx context.Context, assetID string, req IngestRequest) {
    // Update status to ingesting
    s.updateAssetStatus(ctx, assetID, StateIngesting)

    // Create TAMS source
    sourceID, err := s.tams.CreateSource(ctx, TAMSSourceRequest{
        Label:       req.Title,
        Format:      mapTypeToTAMSFormat(req.Type),
        Description: req.Description,
    })
    if err != nil {
        s.handleIngestError(ctx, assetID, "failed to create TAMS source", err)
        return
    }

    // Create TAMS flow
    flowID, err := s.tams.CreateFlow(ctx, TAMSFlowRequest{
        SourceID:    sourceID,
        Label:       req.Title,
        Format:      mapTypeToTAMSFormat(req.Type),
        Description: req.Description,
    })
    if err != nil {
        s.handleIngestError(ctx, assetID, "failed to create TAMS flow", err)
        return
    }

    // Update asset with TAMS IDs
    err = s.db.UpdateAssetTAMSInfo(ctx, assetID, sourceID, flowID)
    if err != nil {
        s.handleIngestError(ctx, assetID, "failed to update asset TAMS info", err)
        return
    }

    // Start media processing
    s.updateAssetStatus(ctx, assetID, StateProcessing)
    err = s.processor.ProcessMedia(ctx, ProcessRequest{
        AssetID:    assetID,
        SourcePath: req.SourcePath,
        FlowID:     flowID,
    })
    if err != nil {
        s.handleIngestError(ctx, assetID, "failed to process media", err)
        return
    }

    // Update status to ready
    s.updateAssetStatus(ctx, assetID, StateReady)
}

// TAMS integration
type TamsClient struct {
    baseURL    string
    httpClient *http.Client
}

func NewTamsClient(baseURL string) *TamsClient {
    return &TamsClient{
        baseURL: baseURL,
        httpClient: &http.Client{
            Timeout: time.Second * 30,
        },
    }
}

func (c *TamsClient) CreateSource(ctx context.Context, req TAMSSourceRequest) (string, error) {
    // Implementation for creating TAMS source
    // ...
    return "source-id", nil
}

func (c *TamsClient) CreateFlow(ctx context.Context, req TAMSFlowRequest) (string, error) {
    // Implementation for creating TAMS flow
    // ...
    return "flow-id", nil
}

// Media processing
type MediaProcessor struct {
    // Add fields for transcoding service, etc.
}

func NewMediaProcessor() *MediaProcessor {
    return &MediaProcessor{}
}

type ProcessRequest struct {
    AssetID    string
    SourcePath string
    FlowID     string
}

func (p *MediaProcessor) ProcessMedia(ctx context.Context, req ProcessRequest) error {
    // Implement media processing:
    // 1. Generate proxies
    // 2. Extract technical metadata
    // 3. Create thumbnails
    // 4. Upload segments to TAMS
    // ...
    return nil
}

// Helper methods
func (s *Service) updateAssetStatus(ctx context.Context, assetID, status string) error {
    err := s.db.UpdateAssetStatus(ctx, assetID, status)
    if err != nil {
        return fmt.Errorf("failed to update asset status: %w", err)
    }
    
    // Invalidate cache
    cacheKey := fmt.Sprintf("asset:%s", assetID)
    s.cache.Del(ctx, cacheKey)
    
    return nil
}

func (s *Service) handleIngestError(ctx context.Context, assetID, message string, err error) {
    // Log error
    log.Printf("Ingest error for asset %s: %s - %v", assetID, message, err)
    
    // Update asset status
    s.updateAssetStatus(ctx, assetID, StateFailed)
    
    // Store error details in asset metadata
    errorInfo := map[string]interface{}{
        "error_message": message,
        "error_detail": err.Error(),
        "error_time":   time.Now(),
    }
    s.db.UpdateAssetErrorInfo(ctx, assetID, errorInfo)
}

func mapTypeToTAMSFormat(assetType string) string {
    switch assetType {
    case "video":
        return "urn:x-nmos:format:video"
    case "audio":
        return "urn:x-nmos:format:audio"
    case "data":
        return "urn:x-nmos:format:data"
    default:
        return "urn:x-nmos:format:data"
    }
}
