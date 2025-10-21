package file

import (
	"context"
	"vcx/agent/internal/infra/db"
	"vcx/agent/internal/infra/db/consts"
	"vcx/agent/internal/infra/db/store"
)

/*
SCHEMA    = { 'ID':           'text PRIMARY KEY',  # Locally created UUID
              'creationDate': 'text NOT NULL',     #

              'path':         'text',              # relative path (NOT Absolute)
              'blobID':       'text',              #
              'branchID':     'text',              # branch pointer at this path
              'changeID':     'text',              # pointer to latest changeID
              'fileID_live':  'text',              # pointer to branch/file that i'm subscribing to
              'permissions':  'text',              # am i allowing this file to be subscribed to (for example)
                                                   # will need to allow individual user-based permissions
                                                   # consider group-based permissions
              'subscribers':  'text',              # who is subscribing to this file
              'isDeleted':    'integer',           # indicates if the file is "available"
                                                   # slight optimization here on restrictions
              'export_restrictions': 'text',       # json of file availability attributes
                                                   # - e.g. available in instance {x} only
              'import_restrictions': 'text'        # json of file merge attributes
                                                   # - e.g. Role Restricted or Branch restricted
              }
*/
const tableName = "file"
const /**Columns*/ (
    COL_ID            = consts.ID
    COL_CREATIONDATE  = consts.CREATIONDATE
    COL_LMU           = consts.LMU
    COL_LMD           = consts.LMD
    COL_GUID          = consts.GUID

    COL_PATH          = consts.PATH
    COL_BLOBID        = consts.BLOBID
    COL_BRANCHID      = consts.BRANCHID
    COL_CHANGEID      = consts.CHANGEID

)


var schema = map[string]string{
    COL_PATH:     consts.TYPE_STRING,
    COL_BLOBID:   consts.TYPE_FOREIGNKEY,
    COL_BRANCHID: consts.TYPE_FOREIGNKEY,
    COL_CHANGEID: consts.TYPE_FOREIGNKEY,
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
