package migrations

import (
	"context"

	"vcx/agent/internal/infra/db"
	fileStore "vcx/agent/internal/infra/db/store/file"
)

func init() {
	Register(Migration{
		Version:     2,
		Description: "Add type and target columns to file table",
		Up: func(ctx context.Context) error {
			if err := db.AddColumn("file", fileStore.COL_TYPE,   "TEXT"); err != nil {
				return err
			}
			if err := db.AddColumn("file", fileStore.COL_TARGET, "TEXT"); err != nil {
				return err
			}
			return nil
		},
		Down: func(ctx context.Context) error {
			// SQLite does not support DROP COLUMN prior to v3.35;
			// no-op here — reset via database file deletion if needed.
			return nil
		},
	})
}
