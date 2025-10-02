package instance

import (
	"context"
	"fmt"
	db "vcx/agent/internal/infra/db/store/instance"
	"vcx/agent/internal/domains"
	"vcx/pkg/toolkit/mapkit"

	"vcx/pkg/logging"
)

var log = logging.GetLogger()
const Domain = "Instance"

type Instance struct {
	domains.Meta
	Path      string
	ProjectID string
	BranchID  string
	ChangeID  string
}

func mapToStruct(data map[string]any) *Instance {
	return &Instance{
		Meta: domains.Meta{
			ID:           mapkit.GetString(data, db.COL_ID),
			CreationDate: mapkit.GetString(data, db.COL_CREATIONDATE),
			LMU:          mapkit.GetString(data, db.COL_LMU),
			LMD:          mapkit.GetString(data, db.COL_LMD),
			GUID:         mapkit.GetString(data, db.COL_GUID),
		},
		Path:      mapkit.GetString(data, db.COL_PATH),
		ProjectID: mapkit.GetString(data, db.COL_PROJECTID),
		BranchID:  mapkit.GetString(data, db.COL_BRANCHID),
		ChangeID:  mapkit.GetString(data, db.COL_CHANGEID),
	}
}

func New(ctx context.Context, path, projectID, branchID, changeID string) (*Instance, error) {
	data := map[string]any{
		db.COL_PATH:      path,
		db.COL_PROJECTID: projectID,
		db.COL_BRANCHID:  branchID,
		db.COL_CHANGEID:  changeID,
	}
	result, err := db.Create(ctx, data)
	if err != nil {
		log.Error(fmt.Sprintf("%s Creation Failed: %v", Domain, err))
		return nil, err
	}

	return mapToStruct(result), nil
}

func (i *Instance) Update(ctx context.Context) error {
	data := map[string]any{
		db.COL_PATH:      i.Path,
		db.COL_PROJECTID: i.ProjectID,
		db.COL_BRANCHID:  i.BranchID,
		db.COL_CHANGEID:  i.ChangeID,
	}
	_, err := db.Update(ctx, i.ID, data)
	if err != nil {
		log.Error(fmt.Sprintf("%s Update Failed: %v", Domain, err))
	}

	return err
}

func GetByID(ctx context.Context, id string) (*Instance, error) {
	data, err := db.GetByID(ctx, id)
	if err != nil {
		log.Error(fmt.Sprintf("%s Retrieval Failed: %v", Domain, err))
		return nil, err
	}

	return mapToStruct(data), nil
}