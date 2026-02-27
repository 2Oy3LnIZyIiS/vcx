// Package blob provides domain models for content-addressable blob storage.
//
// Blobs are immutable content storage that can be shared across multiple files.
// Content is deduplicated by hash - identical content results in the same blob.
//
// Storage strategy:
//   - Small blobs (≤512KB): Stored in database Blob field, FilePath is empty
//   - Large blobs (>512KB): Stored on filesystem, FilePath contains storage location
package blob

import (
	"context"
	"vcx/agent/internal/domains"
	db "vcx/agent/internal/infra/db/store/blob"
	"vcx/pkg/toolkit/mapkit"
)


const Domain = "Blob"


// Blob represents content-addressable storage for file data.
//
// Blobs are immutable and deduplicated by content hash. Multiple files
// can reference the same blob if they have identical content.
//
// Fields:
//   - Blob: Raw content data (empty if stored on filesystem)
//   - FilePath: Storage location on filesystem (empty if stored in DB)
//   - IsCompressed: Whether content is zstd compressed
//   - IsBinary: Whether content is binary (detected by null bytes)
//   - RefCounter: Number of files referencing this blob
type Blob struct {
	domains.Meta
	Blob         []byte
	FilePath     string
	IsCompressed bool
	IsBinary     bool
	RefCounter   int
}

func mapToStruct(data map[string]any) *Blob {
	return &Blob{
		Meta: domains.Meta{
			ID:           mapkit.GetString(data, db.COL_ID),
			CreationDate: mapkit.GetString(data, db.COL_CREATIONDATE),
			LMU:          mapkit.GetString(data, db.COL_LMU),
			LMD:          mapkit.GetString(data, db.COL_LMD),
			GUID:         mapkit.GetString(data, db.COL_GUID),
		},
		Blob:         mapkit.GetBytes(data, db.COL_BLOB),
		FilePath:     mapkit.GetString(data, db.COL_FILEPATH),
		IsCompressed: mapkit.GetBool(data, db.COL_ISCOMPRESSED),
		IsBinary:     mapkit.GetBool(data, db.COL_ISBINARY),
		RefCounter:   mapkit.GetInt(data, db.COL_REFCOUNTER),
	}
}

// New creates a new blob record with a specific ID (content hash).
//
// Parameters:
//   - id: Content hash (SHA256) to use as blob ID
//   - blob: Content data (nil if stored on filesystem)
//   - filePath: Storage path (empty if stored in DB)
//   - isCompressed: Whether content is compressed
//   - isBinary: Whether content is binary
//
// The blob is created with RefCounter=1.
func New(ctx context.Context, id string, blob []byte, filePath string, isCompressed, isBinary bool) (*Blob, error) {
	data := map[string]any{
		db.COL_ID:           id,
		db.COL_BLOB:         blob,
		db.COL_FILEPATH:     filePath,
		db.COL_ISCOMPRESSED: isCompressed,
		db.COL_ISBINARY:     isBinary,
		db.COL_REFCOUNTER:   1,
	}
	result, err := db.Create(ctx, data)
	if err != nil {
		domains.LogError(Domain, "Creation", err)
		return nil, err
	}

	return mapToStruct(result), nil
}

// IncrementRefCounter increments the reference counter for this blob.
// Called when a new file references this blob (deduplication).
func (b *Blob) IncrementRefCounter(ctx context.Context) error {
	b.RefCounter++
	return b.updateRefCounter(ctx)
}

// DecrementRefCounter decrements the reference counter for this blob.
// Called when a file is deleted. If RefCounter reaches 0, blob can be garbage collected.
func (b *Blob) DecrementRefCounter(ctx context.Context) error {
	if b.RefCounter > 0 {
		b.RefCounter--
	}
	return b.updateRefCounter(ctx)
}

func (b *Blob) updateRefCounter(ctx context.Context) error {
	data := map[string]any{
		db.COL_REFCOUNTER: b.RefCounter,
	}
	_, err := db.Update(ctx, b.ID, data)
	if err != nil {
		domains.LogError(Domain, "RefCounter Update", err)
	}
	return err
}

func GetByID(ctx context.Context, id string) (*Blob, error) {
	data, err := db.GetByID(ctx, id)
	if err != nil {
		domains.LogError(Domain, "Retrieval", err)
		return nil, err
	}

	return mapToStruct(data), nil
}
