package tag

import (
	"context"
	"vcx/agent/internal/infra/db"
	"vcx/agent/internal/infra/db/consts"
	"vcx/agent/internal/infra/db/store"
)

/*
SCHEMA    = { 'ID':           'text PRIMARY KEY',  # Locally created UUID
              'creationDate': 'text NOT NULL',     #

              'accountID':     'text',             # account responsible for change
              'name':         'text',              #
              'description':  'text',              #
              'tagType':      'text',              # change for file, branch, or project
              'fileID':       'text',              #
              'branchID':     'text',              #
              'projectID':    'text',              #
              'changeID':     'text'}              #
*/
const tableName = "tag"
const /**Columns*/ (
    COL_ID           = consts.ID
    COL_CREATIONDATE = consts.CREATIONDATE
    COL_LMU          = consts.LMU
    COL_LMD          = consts.LMD
    COL_GUID         = consts.GUID

    COL_ACCOUNTID    = consts.ACCOUNTID
    COL_NAME         = consts.NAME
    COL_DESCRIPTION  = consts.DESCRIPTION
    COL_TAGTYPE      = "tagtype"
    COL_FILEID       = consts.FILEID
    COL_BRANCHID     = consts.BRANCHID
    COL_PROJECTID    = consts.PROJECTID
    COL_CHANGEID     = consts.CHANGEID
)


var schema = map[string]string{
    COL_ACCOUNTID:   consts.TYPE_FOREIGNKEY,
    COL_NAME:        consts.TYPE_STRING,
    COL_DESCRIPTION: consts.TYPE_STRING,
    COL_TAGTYPE:     consts.TYPE_ENUM,
    COL_FILEID:      consts.TYPE_FOREIGNKEY,
    COL_BRANCHID:    consts.TYPE_FOREIGNKEY,
    COL_PROJECTID:   consts.TYPE_FOREIGNKEY,
    COL_CHANGEID:    consts.TYPE_FOREIGNKEY,
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
