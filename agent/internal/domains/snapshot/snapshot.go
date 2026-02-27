package snapshot

import (
	"context"
	"fmt"
	"vcx/agent/internal/domains"
	db "vcx/agent/internal/infra/db/store/snapshot"
	"vcx/pkg/logging"
	"vcx/pkg/toolkit/mapkit"
)

var log = logging.GetLogger()

const Domain = "Snapshot"

type Snapshot struct {
	domains.Meta
	FileID      string
	ChangeID    string
	BlobID      string
	Hash        string
	Size        int
	Summary     string
	LastIndexed string
}

func mapToStruct(data map[string]any) *Snapshot {
	return &Snapshot{
		Meta: domains.Meta{
			ID:           mapkit.GetString(data, db.COL_ID),
			CreationDate: mapkit.GetString(data, db.COL_CREATIONDATE),
			LMU:          mapkit.GetString(data, db.COL_LMU),
			LMD:          mapkit.GetString(data, db.COL_LMD),
			GUID:         mapkit.GetString(data, db.COL_GUID),
		},
		FileID:      mapkit.GetString(data, db.COL_FILEID),
		ChangeID:    mapkit.GetString(data, db.COL_CHANGEID),
		BlobID:      mapkit.GetString(data, db.COL_BLOBID),
		Hash:        mapkit.GetString(data, db.COL_HASH),
		Size:        mapkit.GetInt(data,    db.COL_SIZE),
		Summary:     mapkit.GetString(data, db.COL_SUMMARY),
		LastIndexed: mapkit.GetString(data, db.COL_LASTINDEXED),
	}
}

func New(ctx context.Context, fileID, changeID, blobID, hash string, size int) (*Snapshot, error) {
	data := map[string]any{
		db.COL_FILEID:   fileID,
		db.COL_CHANGEID: changeID,
		db.COL_BLOBID:   blobID,
		db.COL_HASH:     hash,
		db.COL_SIZE:     size,
	}
	result, err := db.Create(ctx, data)
	if err != nil {
		log.Error(fmt.Sprintf("%s Creation Failed: %v", Domain, err))
		return nil, err
	}

	return mapToStruct(result), nil
}

func (s *Snapshot) Update(ctx context.Context) error {
	data := map[string]any{
		db.COL_FILEID:      s.FileID,
		db.COL_CHANGEID:    s.ChangeID,
		db.COL_BLOBID:      s.BlobID,
		db.COL_HASH:        s.Hash,
		db.COL_SIZE:        s.Size,
		db.COL_SUMMARY:     s.Summary,
		db.COL_LASTINDEXED: s.LastIndexed,
	}
	_, err := db.Update(ctx, s.ID, data)
	if err != nil {
		log.Error(fmt.Sprintf("%s Update Failed: %v", Domain, err))
	}

	return err
}

func GetByID(ctx context.Context, id string) (*Snapshot, error) {
	data, err := db.GetByID(ctx, id)
	if err != nil {
		log.Error(fmt.Sprintf("%s Retrieval Failed: %v", Domain, err))
		return nil, err
	}

	return mapToStruct(data), nil
}
