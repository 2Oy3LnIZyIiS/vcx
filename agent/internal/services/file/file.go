package file

import (
	"context"
	"fmt"

	fileDomain "vcx/agent/internal/domains/file"
	blobService "vcx/agent/internal/services/blob"
	tagService "vcx/agent/internal/services/tag"
	"vcx/pkg/logging"
)

var log = logging.GetLogger()

// Ingest creates blob and file record with system tag
func Ingest(ctx context.Context, filePath string) (*fileDomain.File, error) {
	// Create blob (reads file, detects binary, compresses, stores)
	blob, err := blobService.Create(ctx, filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create blob: %w", err)
	}

	// Create file record
	file, err := fileDomain.New(ctx, filePath, blob.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}

	// Create system tag
	_, err = tagService.CreateSystemFileTag(ctx, file.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to create file tag: %w", err)
	}

	log.Debug("Ingested file", "path", filePath, "fileID", file.ID, "blobID", blob.ID)
	return file, nil
}
