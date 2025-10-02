package project

import (
	"context"
	"fmt"
	db "vcx/agent/internal/infra/db/store/project"
	"vcx/agent/internal/domains"
	"vcx/pkg/toolkit/mapkit"

	"vcx/pkg/logging"
)

var log = logging.GetLogger()
const Domain = "Project"


type Project struct {
	domains.Meta  // Embedded metadata fields
	Name     string
	ChangeID string
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
	}
}


func New(ctx context.Context, name, changeID string) (*Project, error) {
	data := map[string]any{
		db.COL_NAME:     name,
		db.COL_CHANGEID: changeID,
	}
    result, err := db.Create(ctx, data)
    if err != nil {
        log.Error(fmt.Sprintf("%s Creation Failed: %v", Domain, err))
        return nil, err
    }

    return mapToStruct(result), nil
}


func (proj *Project) Update(ctx context.Context) error {
	data := map[string]any{
		db.COL_NAME:     proj.Name,
		db.COL_CHANGEID: proj.ChangeID,
	}
    _, err := db.Update(ctx, proj.ID, data)
    if err != nil {
        log.Error(fmt.Sprintf("%s Update Failed: %v", Domain, err))
    }

    return err
}


func GetByID(ctx context.Context, id string) (*Project, error) {
    data, err := db.GetByID(ctx, id)
    if err != nil {
        log.Error(fmt.Sprintf("%s Retrieval Failed: %v", Domain, err))
        return nil, err
    }

    return mapToStruct(data), nil
}
