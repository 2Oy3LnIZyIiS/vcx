package file

import (
	"context"
	"fmt"
	db "vcx/agent/internal/infra/db/store/file"
	"vcx/agent/internal/domains"
	"vcx/pkg/toolkit/mapkit"

	"vcx/pkg/logging"
)

var log = logging.GetLogger()
const Domain = "File"

type File struct {
	domains.Meta
	Path     string
	BlobID   string
	BranchID string
	ChangeID string
}

func mapToStruct(data map[string]any) *File {
	return &File{
		Meta: domains.Meta{
			ID:           mapkit.GetString(data, db.COL_ID),
			CreationDate: mapkit.GetString(data, db.COL_CREATIONDATE),
			LMU:          mapkit.GetString(data, db.COL_LMU),
			LMD:          mapkit.GetString(data, db.COL_LMD),
			GUID:         mapkit.GetString(data, db.COL_GUID),
		},
		Path:     mapkit.GetString(data, db.COL_PATH),
		BlobID:   mapkit.GetString(data, db.COL_BLOBID),
		BranchID: mapkit.GetString(data, db.COL_BRANCHID),
		ChangeID: mapkit.GetString(data, db.COL_CHANGEID),
	}
}

func New(ctx context.Context, path, blobID, branchID, changeID string) (*File, error) {
	data := map[string]any{
		db.COL_PATH:     path,
		db.COL_BLOBID:   blobID,
		db.COL_BRANCHID: branchID,
		db.COL_CHANGEID: changeID,
	}
	result, err := db.Create(ctx, data)
	if err != nil {
		log.Error(fmt.Sprintf("%s Creation Failed: %v", Domain, err))
		return nil, err
	}

	return mapToStruct(result), nil
}

func (f *File) Update(ctx context.Context) error {
	data := map[string]any{
		db.COL_PATH:     f.Path,
		db.COL_BLOBID:   f.BlobID,
		db.COL_BRANCHID: f.BranchID,
		db.COL_CHANGEID: f.ChangeID,
	}
	_, err := db.Update(ctx, f.ID, data)
	if err != nil {
		log.Error(fmt.Sprintf("%s Update Failed: %v", Domain, err))
	}

	return err
}

func GetByID(ctx context.Context, id string) (*File, error) {
	data, err := db.GetByID(ctx, id)
	if err != nil {
		log.Error(fmt.Sprintf("%s Retrieval Failed: %v", Domain, err))
		return nil, err
	}

	return mapToStruct(data), nil
}