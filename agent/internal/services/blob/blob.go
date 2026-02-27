// Package blob provides services for creating and managing content-addressable blobs.
//
// The blob service handles:
//   - Reading file content
//   - Binary detection (null byte check)
//   - Compression (zstd for non-binary files)
//   - Storage strategy (DB vs filesystem based on size)
//   - Deduplication (content-addressable by SHA256)
//   - Filesystem sharding (2-character prefix like Git)
package blob

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	blobDomain "vcx/agent/internal/domains/blob"
	"vcx/agent/internal/infra/db/dbsetup"
	"vcx/pkg/logging"
	"vcx/pkg/toolkit/compressionkit"
	"vcx/pkg/toolkit/cryptokit"
	"vcx/pkg/toolkit/filekit"
)


var log = logging.GetLogger()


const (
	MAX_DB_BLOB_SIZE = 512 * 1024 // 512KB - optimal for cloud sync and SQLite performance
)


// Create reads a file, determines storage strategy, and creates a blob.
//
// Process:
//  1. Read file content and calculate hash (used as blob ID)
//  2. Check if blob exists by ID - if yes, increment RefCounter
//  3. If new: detect binary, compress, store in DB or filesystem
//
// Storage locations:
//   - DB: blob.Blob contains data, blob.FilePath is empty
//   - Filesystem: blob.Blob is nil, blob.FilePath is ~/.vcx/data/blobs/XX/YYYYYY
//
// Returns the blob domain object (existing or newly created).
func Create(ctx context.Context, filePath string) (*blobDomain.Blob, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Calculate hash (this will be the blob ID)
	hashStr := cryptokit.SHA256Hex(data)

	// Check if blob already exists
	existingBlob, err := blobDomain.GetByID(ctx, hashStr)
	if err == nil {
		// Blob exists, increment ref counter
		if err := existingBlob.IncrementRefCounter(ctx); err != nil {
			return nil, fmt.Errorf("failed to increment ref counter: %w", err)
		}
		log.Debug("Reusing existing blob", "hash", hashStr, "refCounter", existingBlob.RefCounter)
		return existingBlob, nil
	}

	isBinary     := filekit.IsBinary(data)
	isCompressed := false

	// Try compression for non-binary data
	if !isBinary {
		data, isCompressed = compressionkit.Compress(data)
	}

	// Decide: DB or filesystem
	if len(data) <= MAX_DB_BLOB_SIZE {
		// Store in DB - no filepath needed
		return blobDomain.New(ctx, hashStr, data, "", isCompressed, isBinary)
	}

	// Store on filesystem - filepath points to storage location
	blobPath, err := writeToDisk(data, hashStr)
	if err != nil {
		return nil, fmt.Errorf("failed to write blob to disk: %w", err)
	}
	return blobDomain.New(ctx, hashStr, nil, blobPath, isCompressed, isBinary)
}


func writeToDisk(data []byte, hashStr string) (string, error) {
	// Shard by first 2 characters (like Git)
	shardDir := filepath.Join(dbsetup.BlobStorePath, hashStr[:2])
	if err := os.MkdirAll(shardDir, 0755); err != nil {
		return "", err
	}

	path := filepath.Join(shardDir, hashStr[2:])
	if err := os.WriteFile(path, data, 0644); err != nil {
		return "", err
	}

	return path, nil
}
