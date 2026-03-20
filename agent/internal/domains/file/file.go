package file

import (
	"context"
	"vcx/agent/internal/consts/filetype"
	"vcx/agent/internal/domains"
	db "vcx/agent/internal/infra/db/store/file"
	"vcx/agent/internal/session"
	"vcx/pkg/toolkit/mapkit"
)

const Domain = "File"

type File struct {
	domains.Meta
	Path      string
	Type      filetype.FileType
	Target    string
	BlobID    string
	BranchID  string
	ChangeID  string
    IsDeleted bool
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
		Path:      mapkit.GetString(data, db.COL_PATH),
		Type:      filetype.FromString(mapkit.GetString(data, db.COL_TYPE)),
		Target:    mapkit.GetString(data, db.COL_TARGET),
		BlobID:    mapkit.GetString(data, db.COL_BLOBID),
		BranchID:  mapkit.GetString(data, db.COL_BRANCHID),
		ChangeID:  mapkit.GetString(data, db.COL_CHANGEID),
		IsDeleted: mapkit.GetBool(data, db.COL_ISDELETED),
	}
}

func New(ctx context.Context, path, blobID string) (*File, error) {
	data := map[string]any{
		db.COL_PATH:      path,
		db.COL_TYPE:      filetype.FILE.ToString(),
		db.COL_BLOBID:    blobID,
		db.COL_BRANCHID:  session.GetBranchID(ctx),
		db.COL_CHANGEID:  session.GetChangeID(ctx),
        db.COL_ISDELETED: false,
	}
	result, err := db.Create(ctx, data)
	if err != nil {
		domains.LogError(Domain, "Creation", err)
		return nil, err
	}

	return mapToStruct(result), nil
}

func NewSymlink(ctx context.Context, path, target string) (*File, error) {
	data := map[string]any{
		db.COL_PATH:      path,
		db.COL_TYPE:      filetype.SYMLINK.ToString(),
		db.COL_TARGET:    target,
		db.COL_BRANCHID:  session.GetBranchID(ctx),
		db.COL_CHANGEID:  session.GetChangeID(ctx),
        db.COL_ISDELETED: false,
	}
	result, err := db.Create(ctx, data)
	if err != nil {
		domains.LogError(Domain, "Creation", err)
		return nil, err
	}

	return mapToStruct(result), nil
}


func (f *File) Update(ctx context.Context) error {
	data := map[string]any{
		db.COL_PATH:      f.Path,
		db.COL_TYPE:      f.Type.ToString(),
		db.COL_TARGET:    f.Target,
		db.COL_BLOBID:    f.BlobID,
		db.COL_BRANCHID:  f.BranchID,
		db.COL_CHANGEID:  f.ChangeID,
		db.COL_ISDELETED: f.IsDeleted,
	}
	_, err := db.Update(ctx, f.ID, data)
	if err != nil {
		domains.LogError(Domain, "Update", err)
	}

	return err
}

func GetByID(ctx context.Context, id string) (*File, error) {
	data, err := db.GetByID(ctx, id)
	if err != nil {
		domains.LogError(Domain, "Retrieval", err)
		return nil, err
	}

	return mapToStruct(data), nil
}
