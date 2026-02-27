package snapshot

import (
	"context"

	"vcx/agent/internal/infra/db"
	"vcx/agent/internal/infra/db/consts"
	"vcx/agent/internal/infra/db/store"
)


const tableName = "snapshot"
const /**Columns*/ (
	COL_ID           = consts.ID
	COL_CREATIONDATE = consts.CREATIONDATE
	COL_LMU          = consts.LMU
	COL_LMD          = consts.LMD
	COL_GUID         = consts.GUID

	COL_FILEID      = "fileID"
	COL_CHANGEID    = "changeID"
	COL_BLOBID      = "blobID"
	COL_HASH        = "hash"
	COL_SIZE        = "size"
	COL_SUMMARY     = "summary"
	COL_LASTINDEXED = "lastIndexed"
)


var schema = map[string]string{
	COL_FILEID:      consts.TYPE_STRING,
	COL_CHANGEID:    consts.TYPE_STRING,
	COL_BLOBID:      consts.TYPE_STRING,
	COL_HASH:        consts.TYPE_STRING,
	COL_SIZE:        consts.TYPE_INT,
	COL_SUMMARY:     consts.TYPE_STRING,
	COL_LASTINDEXED: consts.TYPE_STRING,
}


func CreateTable() {
	db.CreateTable(tableName, schema)
}


func Create(ctx context.Context, data map[string]any) (map[string]any, error) {
	return store.Create(ctx, tableName, data, schema)
}


func Update(ctx context.Context, id string, data map[string]any) (map[string]any, error) {
	return store.Update(ctx, tableName, id, data)
}


func GetByID(ctx context.Context, id string) (map[string]any, error) {
	return store.GetByID(ctx, tableName, id)
}
