package instance

import (
	"context"
	"vcx/agent/internal/infra/db"
	"vcx/agent/internal/infra/db/consts"
	"vcx/agent/internal/infra/db/store"
)

/*
# NOTE: changeID here is usually BLANK.  if a value is set here it is because the
# user has manually reverted/rewound the state of this instance to a particular
# point in time.
# TODO: file changes if a changeID exists for this instance should spawn a new branch
# if changeID is BLANK then both project and branch each have pointers to their
# latest changeIDs
*/
const tableName = "instance"
const /**Columns*/ (
    COL_ID           = consts.ID
    COL_CREATIONDATE = consts.CREATIONDATE
    COL_LMU          = consts.LMU
    COL_LMD          = consts.LMD
    COL_GUID         = consts.GUID

    COL_PATH         = consts.PATH
    COL_PROJECTID    = consts.PROJECTID
    COL_BRANCHID     = consts.BRANCHID
    COL_CHANGEID     = consts.CHANGEID
)


var schema = map[string]string{
    COL_PATH:      consts.TYPE_STRING,
    COL_PROJECTID: consts.TYPE_FOREIGNKEY,
    COL_BRANCHID:  consts.TYPE_FOREIGNKEY,
    COL_CHANGEID:  consts.TYPE_FOREIGNKEY,
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
