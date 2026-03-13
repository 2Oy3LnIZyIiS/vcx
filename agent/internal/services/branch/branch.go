package branch

import (
	"context"

	branchDomain "vcx/agent/internal/domains/branch"
	"vcx/pkg/logging"
)

var log = logging.GetLogger()


func Create(ctx context.Context, name string) (*branchDomain.Branch, error) {
	return branchDomain.New(ctx, name)
}


func GetByID(ctx context.Context, id string) (*branchDomain.Branch, error) {
	return branchDomain.GetByID(ctx, id)
}
