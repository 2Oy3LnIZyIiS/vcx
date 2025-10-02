package change

import (
	"context"
	"vcx/agent/internal/infra/db"
	"vcx/agent/internal/infra/db/consts"
	"vcx/agent/internal/infra/db/store"
)

const tableName = "change"
const /**Columns*/ (
    COL_ID            = consts.ID
    COL_CREATIONDATE  = consts.CREATIONDATE
    COL_LMU           = consts.LMU
    COL_LMD           = consts.LMD
    COL_GUID          = consts.GUID

    COL_ACCOUNTID     = consts.ACCOUNTID
    COL_FILEID        = consts.FILEID
    COL_BRANCHID      = consts.BRANCHID
    COL_PROJECTID     = consts.PROJECTID
    COL_CHANGETYPE    = "changeType"
    COL_CHANGEBLOB    = "changeBlob"
    COL_CHANGEID_PREV = "changeID_prev"
    COL_CHANGEID_NEXT = "changeID_next"
)


var schema = map[string]string{
    COL_ACCOUNTID:     consts.TYPE_FOREIGNKEY,
    COL_FILEID:        consts.TYPE_FOREIGNKEY,
    COL_BRANCHID:      consts.TYPE_FOREIGNKEY,
    COL_PROJECTID:     consts.TYPE_FOREIGNKEY,
    COL_CHANGETYPE:    consts.TYPE_ENUM,
    COL_CHANGEBLOB:    consts.TYPE_BLOB,
    COL_CHANGEID_PREV: consts.TYPE_FOREIGNKEY,
    COL_CHANGEID_NEXT: consts.TYPE_FOREIGNKEY,
}


func init()  {
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
