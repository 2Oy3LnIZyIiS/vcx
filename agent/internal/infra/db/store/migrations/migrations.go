package migrations

import (
	"context"

	"vcx/agent/internal/infra/db"
	"vcx/agent/internal/infra/db/consts"
	"vcx/agent/internal/infra/db/store"
)

const tableName = "migrations"

const /**Columns*/ (
	COL_ID           = consts.ID
	COL_CREATIONDATE = consts.CREATIONDATE
	COL_LMU          = consts.LMU
	COL_LMD          = consts.LMD
	COL_GUID         = consts.GUID

	COL_VERSION      = "version"
	COL_APPLIED_AT   = "applied_at"
)

var schema = map[string]string{
	COL_VERSION:    "INTEGER NOT NULL",
	COL_APPLIED_AT: "DATETIME DEFAULT CURRENT_TIMESTAMP",
}

func CreateTable() {
	db.CreateTable(tableName, schema)
}

func Create(ctx context.Context, data map[string]any) (map[string]any, error) {
	return store.Create(ctx, tableName, data, schema)
}

func GetMaxVersion(ctx context.Context) (int, error) {
	// Try to get max version using a simple select
	results, err := db.SelectWithContext(ctx, tableName, []string{"COALESCE(MAX(version), 0) as max_version"}, map[string]any{})
	if err != nil || len(results) == 0 {
		return 0, nil // Table empty or doesn't exist
	}
	
	if maxVersion, ok := results[0]["max_version"]; ok && maxVersion != nil {
		switch v := maxVersion.(type) {
		case int:
			return v, nil
		case int64:
			return int(v), nil
		case float64:
			return int(v), nil
		}
	}
	return 0, nil
}