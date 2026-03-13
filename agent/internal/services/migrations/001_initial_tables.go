package migrations

import (
	"context"

	"vcx/agent/internal/infra/db/store/account"
	"vcx/agent/internal/infra/db/store/blob"
	"vcx/agent/internal/infra/db/store/branch"
	"vcx/agent/internal/infra/db/store/change"
	"vcx/agent/internal/infra/db/store/file"
	"vcx/agent/internal/infra/db/store/instance"
	"vcx/agent/internal/infra/db/store/project"
	"vcx/agent/internal/infra/db/store/simplekv"
	"vcx/agent/internal/infra/db/store/tag"
)

func init() {
	Register(Migration{
		Version:     1,
		Description: "Create initial tables",
		Up: func(ctx context.Context) error {
			account.CreateTable()
			blob.CreateTable()
			branch.CreateTable()
			change.CreateTable()
			file.CreateTable()
			instance.CreateTable()
			project.CreateTable()
			simplekv.CreateTable()
			tag.CreateTable()
			return nil
		},
		Down: func(ctx context.Context) error {
			// Drop tables in reverse order
			tables := []string{
				"tag", "simplekv", "project", "instance",
				"file", "change", "branch", "blob", "account",
			}
			for _, table := range tables {
				log.Debug("Would drop table", "table", table)
				// Note: Can't actually drop tables with current db package
				// In practice, you'd delete the database file for reset
			}
			return nil
		},
	})
}
