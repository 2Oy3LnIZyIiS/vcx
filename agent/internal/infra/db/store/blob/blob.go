package blob

import (
	"context"

	"vcx/agent/internal/infra/db"
	"vcx/agent/internal/infra/db/consts"
	"vcx/agent/internal/infra/db/store"
)


const tableName = "blob"
const /**Columns*/ (
    COL_ID           = consts.ID
    COL_CREATIONDATE = consts.CREATIONDATE
    COL_LMU          = consts.LMU
    COL_LMD          = consts.LMD
    COL_GUID         = consts.GUID

    COL_BLOB         = "blob"
    COL_FILEPATH     = "filepath"
    COL_ISCOMPRESSED = "isCompressed"
    COL_ISBINARY     = "isBinary"
    COL_REFCOUNTER   = "refCounter"
)


var schema = map[string]string{
	COL_BLOB:         consts.TYPE_BLOB,
    COL_FILEPATH:     consts.TYPE_STRING,
	COL_ISCOMPRESSED: consts.TYPE_BOOL,
	COL_ISBINARY:     consts.TYPE_BOOL,
	COL_REFCOUNTER:   "integer",
}


func CreateTable()  {
    db.CreateTable(tableName, schema)
}


func Create(ctx context.Context, data map[string]any,) (map[string]any, error) {
    return store.Create(ctx, tableName, data, schema)
}


func Update(ctx context.Context, id string, data map[string]any) (map[string]any, error) {
    return store.Update(ctx, tableName, id, data)
}


func GetByID(ctx context.Context, id string) (map[string]any, error) {
    return store.GetByID(ctx, tableName, id)
}
