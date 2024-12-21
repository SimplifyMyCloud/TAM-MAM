// processor/media_processor.go
package processor

import (
    "context"
    "encoding/json"
    "fmt"
    "os"
    "path/filepath"
    "time"

    "github.com/google/uuid"
)

type MediaProcessor struct {
    config        ProcessorConfig
    transcoder    *Transcoder
    metadataExtractor *MetadataExtractor
    tamsClient    *TamsClient
    storage       *Storage
}

type ProcessorConfig struct {
    WorkDir           string
    ProxyFormats     []ProxyFormat
    SegmentDuration  int // in seconds
    MaxConcurrent    int
    FFmpegPath       string
}

type ProxyFormat struct {
    Name     string
    Width    int
    Height   int
    Bitrate  string
    Codec    string
}

func NewMediaProcessor(config ProcessorConfig, tamsClient *TamsClient) *MediaProcessor {
    return &MediaProcessor{
        config:     config,
        transcoder: NewTranscoder(config.FFmpegPath),
        metadataExtractor: NewMetadataExtractor(),
        tamsClient: tamsClient,
        storage:   NewStorage(),
    }
}

func (p *MediaProcessor) ProcessMedia(ctx context.Context, req ProcessRequest) error {
    // Create working directory for this asset
    workDir := filepath.Join(p.config.WorkDir, req.AssetID)
    if err := os.MkdirAll(workDir, 0755); err != nil {
        return fmt.Errorf("failed to create work directory: %w", err)
    }
    defer os.RemoveAll(workDir)

    // Extract technical metadata
    metadata, err := p.metadataExtractor.Extract(req.SourcePath)
    if err != nil {
        return fmt.Errorf("failed to extract metadata: %w", err)
    }

    // Generate thumbnails
    thumbnailPath, err := p.generateThumbnails(req.SourcePath, workDir)
    if err != nil {
        return fmt.Errorf("failed to generate thumbnails: %w", err)
    }

    // Create proxies
    proxyPaths, err := p.createProxies(ctx, req.SourcePath, workDir)
    if err != nil {
        return fmt.Errorf("failed to create proxies: %w", err)
    }

    // Segment and upload to TAMS
    err = p.segmentAndUpload(ctx, req.SourcePath, req.FlowID, metadata)
    if err != nil {
        return fmt.Errorf("failed to segment and upload: %w", err)
    }

    return nil
}

type Transcoder struct {
    ffmpegPath string
}

func NewTranscoder(ffmpegPath string) *Transcoder {
    return &Transcoder{ffmpegPath: ffmpegPath}
}

func (t *Transcoder) CreateProxy(ctx context.Context, input string, output string, format ProxyFormat) error {
    args := []string{
        "-i", input,
        "-c:v", format.Codec,
        "-b:v", format.Bitrate,
        "-vf", fmt.Sprintf("scale=%d:%d", format.Width, format.Height),
        "-y",
        output,
    }

    cmd := exec.CommandContext(ctx, t.ffmpegPath, args...)
    return cmd.Run()
}

func (t *Transcoder) Segment(ctx context.Context, input string, segmentPattern string, duration int) error {
    args := []string{
        "-i", input,
        "-c", "copy",
        "-f", "segment",
        "-segment_time", fmt.Sprintf("%d", duration),
        "-reset_timestamps", "1",
        segmentPattern,
    }

    cmd := exec.CommandContext(ctx, t.ffmpegPath, args...)
    return cmd.Run()
}

type MetadataExtractor struct{}

func NewMetadataExtractor() *MetadataExtractor {
    return &MetadataExtractor{}
}

type MediaMetadata struct {
    Duration    float64
    Format      string
    Codec       string
    Width       int
    Height      int
    Framerate   string
    Bitrate     int64
    AudioCodec  string
    AudioRate   int
    AudioLayout string
}

func (e *MetadataExtractor) Extract(filepath string) (*MediaMetadata, error) {
    args := []string{
        "-v", "quiet",
        "-print_format", "json",
        "-show_format",
        "-show_streams",
        filepath,
    }

    cmd := exec.Command("ffprobe", args...)
    output, err := cmd.Output()
    if err != nil {
        return nil, fmt.Errorf("ffprobe failed: %w", err)
    }

    var ffprobeData map[string]interface{}
    if err := json.Unmarshal(output, &ffprobeData); err != nil {
        return nil, fmt.Errorf("failed to parse ffprobe output: %w", err)
    }

    // Parse ffprobe output into MediaMetadata
    metadata := &MediaMetadata{
        // Fill in metadata from ffprobeData
    }

    return metadata, nil
}

func (p *MediaProcessor) generateThumbnails(input string, workDir string) (string, error) {
    thumbnailPath := filepath.Join(workDir, "thumbnail.jpg")
    args := []string{
        "-i", input,
        "-vf", "thumbnail,scale=640:360",
        "-frames:v", "1",
        thumbnailPath,
    }

    cmd := exec.Command(p.config.FFmpegPath, args...)
    if err := cmd.Run(); err != nil {
        return "", fmt.Errorf("failed to generate thumbnail: %w", err)
    }

    return thumbnailPath, nil
}

func (p *MediaProcessor) createProxies(ctx context.Context, input string, workDir string) (map[string]string, error) {
    proxyPaths := make(map[string]string)
    for _, format := range p.config.ProxyFormats {
        outputPath := filepath.Join(workDir, fmt.Sprintf("proxy_%s.mp4", format.Name))
        err := p.transcoder.CreateProxy(ctx, input, outputPath, format)
        if err != nil {
            return nil, fmt.Errorf("failed to create %s proxy: %w", format.Name, err)
        }
        proxyPaths[format.Name] = outputPath
    }
    return proxyPaths, nil
}

func (p *MediaProcessor) segmentAndUpload(ctx context.Context, input string, flowID string, metadata *MediaMetadata) error {
    // Create temporary directory for segments
    segmentDir := filepath.Join(p.config.WorkDir, "segments")
    if err := os.MkdirAll(segmentDir, 0755); err != nil {
        return fmt.Errorf("failed to create segment directory: %w", err)
    }
    defer os.RemoveAll(segmentDir)

    // Segment the file
    segmentPattern := filepath.Join(segmentDir, "segment_%04d.ts")
    err := p.transcoder.Segment(ctx, input, segmentPattern, p.config.SegmentDuration)
    if err != nil {
        return fmt.Errorf("failed to segment file: %w", err)
    }

    // Upload segments to TAMS
    segments, err := filepath.Glob(filepath.Join(segmentDir, "segment_*.ts"))
    if err != nil {
        return fmt.Errorf("failed to list segments: %w", err)
    }

    for i, segment := range segments {
        // Get segment storage location from TAMS
        storageInfo, err := p.tamsClient.AllocateStorage(ctx, flowID)
        if err != nil {
            return fmt.Errorf("failed to allocate storage: %w", err)
        }

        // Upload segment
        if err := p.storage.UploadSegment(ctx, segment, storageInfo.PutURL); err != nil {
            return fmt.Errorf("failed to upload segment: %w", err)
        }

        // Register segment with TAMS
        startTime := time.Duration(i) * time.Second * time.Duration(p.config.SegmentDuration)
        endTime := time.Duration(i+1) * time.Second * time.Duration(p.config.SegmentDuration)
        
        err = p.tamsClient.RegisterSegment(ctx, flowID, TAMSSegment{
            ObjectID:  storageInfo.ObjectID,
            TimeRange: fmt.Sprintf("%v:%v", startTime, endTime),
        })
        if err != nil {
            return fmt.Errorf("failed to register segment: %w", err)
        }
    }

    return nil
}
