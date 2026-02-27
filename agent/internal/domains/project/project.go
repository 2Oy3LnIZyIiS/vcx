package project

import (
	"context"
	db "vcx/agent/internal/infra/db/store/project"
	"vcx/agent/internal/domains"
	"vcx/agent/internal/session"
	"vcx/pkg/toolkit/mapkit"
)

const Domain = "Project"


type Project struct {
	domains.Meta  // Embedded metadata fields
	Name     string
	ChangeID string
	DefaultBranchID string
}


func mapToStruct(data map[string]any) *Project {
	return &Project{
		Meta: domains.Meta{
			ID:           mapkit.GetString(data, db.COL_ID),
			CreationDate: mapkit.GetString(data, db.COL_CREATIONDATE),
			LMU:          mapkit.GetString(data, db.COL_LMU),
			LMD:          mapkit.GetString(data, db.COL_LMD),
			GUID:         mapkit.GetString(data, db.COL_GUID),
		},
		Name:     mapkit.GetString(data, db.COL_NAME),
		ChangeID: mapkit.GetString(data, db.COL_CHANGEID),
		DefaultBranchID: mapkit.GetString(data, db.COL_DEFAULTBRANCHID),
	}
}


func New(ctx context.Context, name string) (*Project, error) {
	data := map[string]any{
		db.COL_NAME:     name,
		db.COL_CHANGEID: session.GetChangeID(ctx),
		db.COL_DEFAULTBRANCHID: "",
	}
    result, err := db.Create(ctx, data)
    if err != nil {
        domains.LogError(Domain, "Creation", err)
        return nil, err
    }

    return mapToStruct(result), nil
}


func (proj *Project) Update(ctx context.Context) error {
	data := map[string]any{
		db.COL_NAME:     proj.Name,
		db.COL_CHANGEID: proj.ChangeID,
		db.COL_DEFAULTBRANCHID: proj.DefaultBranchID,
	}
    _, err := db.Update(ctx, proj.ID, data)
    if err != nil {
        domains.LogError(Domain, "Update", err)
    }

    return err
}


func GetByID(ctx context.Context, id string) (*Project, error) {
    data, err := db.GetByID(ctx, id)
    if err != nil {
        domains.LogError(Domain, "Retrieval", err)
        return nil, err
    }

    return mapToStruct(data), nil
}
