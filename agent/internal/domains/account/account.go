package account

import (
	"context"
	"fmt"
	db "vcx/agent/internal/infra/db/store/account"
	"vcx/agent/internal/domains"
	"vcx/pkg/toolkit/mapkit"

	"vcx/pkg/logging"
)

var log = logging.GetLogger()
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
        log.Error(fmt.Sprintf("%s Creation Failed: %v", Domain, err))
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
        log.Error(fmt.Sprintf("%s Update Failed: %v", Domain, err))
    }

    return err
}


func GetByID(ctx context.Context, id string) (*Account, error) {
    data, err := db.GetByID(ctx, id)
    if err != nil {
        log.Error(fmt.Sprintf("%s Retrieval Failed: %v", Domain, err))
        return nil, err
    }

    return mapToStruct(data), nil
}
