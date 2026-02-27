package account

import (
	"context"
	"vcx/agent/internal/domains"
	db "vcx/agent/internal/infra/db/store/account"
	"vcx/pkg/toolkit/mapkit"
)

const Domain = "Account"


type Account struct {
	domains.Meta  // Embedded metadata fields
	Name    string
	Alias   string
	Email   string
	Display string
}


func mapToStruct(data map[string]any) *Account {
	return &Account{
		Meta: domains.Meta{
			ID:           mapkit.GetString(data, db.COL_ID),
			CreationDate: mapkit.GetString(data, db.COL_CREATIONDATE),
			LMU:          mapkit.GetString(data, db.COL_LMU),
			LMD:          mapkit.GetString(data, db.COL_LMD),
			GUID:         mapkit.GetString(data, db.COL_GUID),
		},
		Name:    mapkit.GetString(data, db.COL_NAME),
		Alias:   mapkit.GetString(data, db.COL_ALIAS),
		Email:   mapkit.GetString(data, db.COL_EMAIL),
		Display: mapkit.GetString(data, db.COL_DISPLAY),
	}
}


func New(ctx context.Context, name, email, alias string) (*Account, error) {
	data := map[string]any{
		db.COL_NAME:  name,
		db.COL_EMAIL: email,
		db.COL_ALIAS: alias,
	}
    result, err := db.Create(ctx, data)
    if err != nil {
        domains.LogError(Domain, "Creation", err)
        return nil, err
    }

    return mapToStruct(result), nil
}


func (acc *Account) Update(ctx context.Context) error {
	data := map[string]any{
		db.COL_NAME:  acc.Name,
		db.COL_EMAIL: acc.Email,
		db.COL_ALIAS: acc.Alias,
	}
    _, err := db.Update(ctx, acc.ID, data)
    if err != nil {
        domains.LogError(Domain, "Update", err)
    }

    return err
}


func GetByID(ctx context.Context, id string) (*Account, error) {
    data, err := db.GetByID(ctx, id)
    if err != nil {
        domains.LogError(Domain, "Retrieval", err)
        return nil, err
    }

    return mapToStruct(data), nil
}
