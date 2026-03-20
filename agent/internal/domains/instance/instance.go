package instance

import (
	"context"
	"vcx/agent/internal/domains"
	db "vcx/agent/internal/infra/db/store/instance"
	"vcx/agent/internal/session"
	"vcx/pkg/toolkit/mapkit"
)

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


func New(ctx context.Context, path string) (*Instance, error) {
	data := map[string]any{
		db.COL_PATH:      path,
		db.COL_PROJECTID: session.GetProjectID(ctx),
		db.COL_BRANCHID:  session.GetBranchID(ctx),
		db.COL_CHANGEID:  session.GetChangeID(ctx),
	}
	result, err := db.Create(ctx, data)
	if err != nil {
		domains.LogError(Domain, "Creation", err)
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
		domains.LogError(Domain, "Update", err)
	}

	return err
}


func GetByID(ctx context.Context, id string) (*Instance, error) {
	data, err := db.GetByID(ctx, id)
	if err != nil {
		domains.LogError(Domain, "Retrieval", err)
		return nil, err
	}

	return mapToStruct(data), nil
}