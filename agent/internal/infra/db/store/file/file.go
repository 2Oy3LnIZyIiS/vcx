package file

import (
	"context"
	"vcx/agent/internal/infra/db"
	"vcx/agent/internal/infra/db/consts"
	"vcx/agent/internal/infra/db/store"
)

const tableName = "file"
const /*Columns*/ (
    COL_ID           = consts.ID
    COL_CREATIONDATE = consts.CREATIONDATE
    COL_LMU          = consts.LMU
    COL_LMD          = consts.LMD
    COL_GUID         = consts.GUID

    COL_PATH         = consts.PATH
    COL_BLOBID       = consts.BLOBID
    COL_BRANCHID     = consts.BRANCHID
    COL_CHANGEID     = consts.CHANGEID
    COL_ISDELETED    = "isDeleted"
)


var schema = map[string]string{
    COL_PATH:      consts.TYPE_STRING,
    COL_BLOBID:    consts.TYPE_FOREIGNKEY,
    COL_BRANCHID:  consts.TYPE_FOREIGNKEY,
    COL_CHANGEID:  consts.TYPE_FOREIGNKEY,
    COL_ISDELETED: consts.TYPE_BOOL,
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
