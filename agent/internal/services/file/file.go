package file

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	fileDomain "vcx/agent/internal/domains/file"
	blobService "vcx/agent/internal/services/blob"
	tagService "vcx/agent/internal/services/tag"
	"vcx/pkg/logging"
)

var log = logging.GetLogger()

// Ingest creates blob and file record with system tag
func Ingest(ctx context.Context, projectPath, filePath string) (*fileDomain.File, error) {
	// Create blob (reads file, detects binary, compresses, stores)
	blob, err := blobService.Create(ctx, filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create blob: %w", err)
	}

	// Store path relative to the project root
	relPath, err := filepath.Rel(projectPath, filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to compute relative path: %w", err)
	}

	// Create file record
	file, err := fileDomain.New(ctx, relPath, blob.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}

	// Create system tag
	_, err = tagService.CreateSystemFileTag(ctx, file.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to create file tag: %w", err)
	}

	log.Debug("Ingested file", "path", relPath, "fileID", file.ID, "blobID", blob.ID)
	return file, nil
}


// IngestSymlink records a symlink's path and target without hashing content.
// The target is stored as-is (absolute or relative) since it may point outside the project.
func IngestSymlink(ctx context.Context, projectPath, linkPath string) (*fileDomain.File, error) {
	target, err := os.Readlink(linkPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read symlink target: %w", err)
	}

	relPath, err := filepath.Rel(projectPath, linkPath)
	if err != nil {
		return nil, fmt.Errorf("failed to compute relative path: %w", err)
	}

	file, err := fileDomain.NewSymlink(ctx, relPath, target)
	if err != nil {
		return nil, fmt.Errorf("failed to create symlink record: %w", err)
	}

	log.Debug("Ingested symlink", "path", relPath, "target", target, "fileID", file.ID)
	return file, nil
}
