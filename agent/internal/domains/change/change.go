package change

import (
	"context"
	"fmt"
	db "vcx/agent/internal/infra/db/store/change"
	"vcx/agent/internal/domains"
	"vcx/pkg/toolkit/mapkit"

	"vcx/pkg/logging"
)

var log = logging.GetLogger()
const Domain = "Change"

type Change struct {
	domains.Meta
	AccountID     string
	FileID        string
	BranchID      string
	ProjectID     string
	ChangeType    string
	ChangeBlob    []byte
	ChangeIDPrev  string
	ChangeIDNext  string
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
		ChangeType:   mapkit.GetString(data, db.COL_CHANGETYPE),
		ChangeBlob:   mapkit.GetBytes(data, db.COL_CHANGEBLOB),
		ChangeIDPrev: mapkit.GetString(data, db.COL_CHANGEID_PREV),
		ChangeIDNext: mapkit.GetString(data, db.COL_CHANGEID_NEXT),
	}
}

func New(ctx context.Context, accountID, fileID, branchID, projectID, changeType string, changeBlob []byte) (*Change, error) {
	data := map[string]any{
		db.COL_ACCOUNTID:  accountID,
		db.COL_FILEID:     fileID,
		db.COL_BRANCHID:   branchID,
		db.COL_PROJECTID:  projectID,
		db.COL_CHANGETYPE: changeType,
		db.COL_CHANGEBLOB: changeBlob,
	}
	result, err := db.Create(ctx, data)
	if err != nil {
		log.Error(fmt.Sprintf("%s Creation Failed: %v", Domain, err))
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
		db.COL_CHANGETYPE:    c.ChangeType,
		db.COL_CHANGEBLOB:    c.ChangeBlob,
		db.COL_CHANGEID_PREV: c.ChangeIDPrev,
		db.COL_CHANGEID_NEXT: c.ChangeIDNext,
	}
	_, err := db.Update(ctx, c.ID, data)
	if err != nil {
		log.Error(fmt.Sprintf("%s Update Failed: %v", Domain, err))
	}

	return err
}

func GetByID(ctx context.Context, id string) (*Change, error) {
	data, err := db.GetByID(ctx, id)
	if err != nil {
		log.Error(fmt.Sprintf("%s Retrieval Failed: %v", Domain, err))
		return nil, err
	}

	return mapToStruct(data), nil
}