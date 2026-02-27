package branch

import (
	"context"
	"vcx/agent/internal/domains"
	"vcx/agent/internal/session"
	db "vcx/agent/internal/infra/db/store/branch"
	"vcx/pkg/toolkit/mapkit"
)

const Domain = "Branch"

type Branch struct {
	domains.Meta
	Name      string
	ProjectID string
	ChangeID  string
}

func mapToStruct(data map[string]any) *Branch {
	return &Branch{
		Meta: domains.Meta{
			ID:           mapkit.GetString(data, db.COL_ID),
			CreationDate: mapkit.GetString(data, db.COL_CREATIONDATE),
			LMU:          mapkit.GetString(data, db.COL_LMU),
			LMD:          mapkit.GetString(data, db.COL_LMD),
			GUID:         mapkit.GetString(data, db.COL_GUID),
		},
		Name:      mapkit.GetString(data, db.COL_NAME),
		ProjectID: mapkit.GetString(data, db.COL_PROJECTID),
		ChangeID:  mapkit.GetString(data, db.COL_CHANGEID),
	}
}

func New(ctx context.Context, name string) (*Branch, error) {
	data := map[string]any{
		db.COL_NAME:      name,
		db.COL_PROJECTID: session.GetProjectID(ctx),
		db.COL_CHANGEID:  session.GetChangeID(ctx),
	}
	result, err := db.Create(ctx, data)
	if err != nil {
		domains.LogError(Domain, "Creation", err)
		return nil, err
	}

	return mapToStruct(result), nil
}

func (b *Branch) Update(ctx context.Context) error {
	data := map[string]any{
		db.COL_NAME:      b.Name,
		db.COL_PROJECTID: b.ProjectID,
		db.COL_CHANGEID:  b.ChangeID,
	}
	_, err := db.Update(ctx, b.ID, data)
	if err != nil {
		domains.LogError(Domain, "Update", err)
	}

	return err
}

func GetByID(ctx context.Context, id string) (*Branch, error) {
	data, err := db.GetByID(ctx, id)
	if err != nil {
		domains.LogError(Domain, "Retrieval", err)
		return nil, err
	}

	return mapToStruct(data), nil
}
