package account

import (
	"context"

	"vcx/agent/internal/infra/db"
	"vcx/agent/internal/infra/db/consts"
	"vcx/agent/internal/infra/db/store"
)


const tableName = "account"
const /**Columns*/ (
    COL_ID           = consts.ID
    COL_CREATIONDATE = consts.CREATIONDATE
    COL_LMU          = consts.LMU
    COL_LMD          = consts.LMD
    COL_GUID         = consts.GUID

    COL_NAME         = consts.NAME
    COL_ALIAS        = "alias"
    COL_EMAIL        = "email"
    COL_DISPLAY      = "display"
)


var schema = map[string]string{
	COL_NAME:         consts.TYPE_STRING,
	COL_ALIAS:        consts.TYPE_STRING,
	COL_EMAIL:        consts.TYPE_STRING,
    COL_DISPLAY:      consts.TYPE_STRING,
}

func CreateTable() {
    db.CreateTable(tableName, schema)
}



func Create(ctx context.Context, data map[string]any,) (map[string]any, error) {
    return store.Create(ctx, tableName, data, schema)
	// id, err := db.InsertWithContext(ctx, tableName, data, schema)
    // if err != nil {
    //     return nil, err
    // }
    // return GetByID(ctx, id)
}


func Update(ctx context.Context, id string, data map[string]any) (map[string]any, error) {
    return store.Update(ctx, tableName, id, data)
	// _, err := db.UpdateWithContext(ctx, tableName, data, map[string]any{consts.ID: id})
    // if err != nil {
    //     return nil, err
    // }
	// return GetByID(ctx, id)
}


func GetByID(ctx context.Context, id string) (map[string]any, error) {
    return store.GetByID(ctx, tableName, id)
	// result, err := db.SelectOneWithContext(ctx, tableName, []string{"*"}, map[string]any{consts.ID: id})
	// if err != nil {
	// 	return nil, err
	// }
	// return result, nil
}
