package change

import (
	"context"

	changeDomain "vcx/agent/internal/domains/change"
	"vcx/pkg/logging"
)

var log = logging.GetLogger()


func CreateProjectChange(ctx context.Context) (*changeDomain.Change, error) {
	return changeDomain.NewProject(ctx)
}


func GetByID(ctx context.Context, id string) (*changeDomain.Change, error) {
	return changeDomain.GetByID(ctx, id)
}
