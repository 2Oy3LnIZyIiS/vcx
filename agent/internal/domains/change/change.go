package change

import (
	"context"
	"vcx/agent/internal/consts/changetype"
	"vcx/agent/internal/domains"
	db "vcx/agent/internal/infra/db/store/change"
	"vcx/agent/internal/session"
	"vcx/pkg/toolkit/mapkit"
)

const Domain = "Change"

type Change struct {
	domains.Meta
	AccountID     string
	FileID        string
	BranchID      string
	ProjectID     string
	ChangeType    changetype.ChangeType
	ChangeBlob    []byte
	ChangeIDPrev  string
	ChangeIDNext  string
	Summary       string
	Embedding     []byte
}

func mapToStruct(data map[string]any) *Change {
	return &Change{
		Meta: domains.Meta{
			ID:           mapkit.GetString(data, db.COL_ID),
			CreationDate: mapkit.GetString(data, db.COL_CREATIONDATE),
			LMU:          mapkit.GetString(data, db.COL_LMU),
			LMD:          mapkit.GetString(data, db.COL_LMD),
			GUID:         mapkit.GetString(data, db.COL_GUID),
		},
		AccountID:    mapkit.GetString(data, db.COL_ACCOUNTID),
		FileID:       mapkit.GetString(data, db.COL_FILEID),
		BranchID:     mapkit.GetString(data, db.COL_BRANCHID),
		ProjectID:    mapkit.GetString(data, db.COL_PROJECTID),
		ChangeType:   changetype.FromString( mapkit.GetString(data, db.COL_CHANGETYPE) ),
		ChangeBlob:   mapkit.GetBytes(data,  db.COL_CHANGEBLOB),
		ChangeIDPrev: mapkit.GetString(data, db.COL_CHANGEID_PREV),
		ChangeIDNext: mapkit.GetString(data, db.COL_CHANGEID_NEXT),
		Summary:      mapkit.GetString(data, db.COL_SUMMARY),
		Embedding:    mapkit.GetBytes(data,  db.COL_EMBEDDING),
	}
}


func NewProject(ctx context.Context) (*Change, error) {
    return New(ctx, session.GetAccountID(ctx), "", "", "", changetype.ACCOUNT, nil)
}


func New(ctx context.Context, accountID, fileID, branchID, projectID string, _changeType changetype.ChangeType, changeBlob []byte) (*Change, error) {
	data := map[string]any{
		db.COL_ACCOUNTID:  accountID,
		db.COL_FILEID:     fileID,
		db.COL_BRANCHID:   branchID,
		db.COL_PROJECTID:  projectID,
		db.COL_CHANGETYPE: _changeType.ToString(),
		db.COL_CHANGEBLOB: changeBlob,
	}
	result, err := db.Create(ctx, data)
	if err != nil {
		domains.LogError(Domain, "Creation", err)
		return nil, err
	}

	return mapToStruct(result), nil
}


func (c *Change) Update(ctx context.Context) error {
	data := map[string]any{
		db.COL_ACCOUNTID:     c.AccountID,
		db.COL_FILEID:        c.FileID,
		db.COL_BRANCHID:      c.BranchID,
		db.COL_PROJECTID:     c.ProjectID,
		db.COL_CHANGETYPE:    c.ChangeType.ToString(),
		db.COL_CHANGEBLOB:    c.ChangeBlob,
		db.COL_CHANGEID_PREV: c.ChangeIDPrev,
		db.COL_CHANGEID_NEXT: c.ChangeIDNext,
		db.COL_SUMMARY:       c.Summary,
		db.COL_EMBEDDING:     c.Embedding,
	}
	_, err := db.Update(ctx, c.ID, data)
	if err != nil {
		domains.LogError(Domain, "Update", err)
	}

	return err
}

func GetByID(ctx context.Context, id string) (*Change, error) {
	data, err := db.GetByID(ctx, id)
	if err != nil {
		domains.LogError(Domain, "Retrieval", err)
		return nil, err
	}

	return mapToStruct(data), nil
}
