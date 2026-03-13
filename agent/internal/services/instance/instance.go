package instance

import (
	"context"

	instanceDomain "vcx/agent/internal/domains/instance"
	"vcx/pkg/logging"
)

var log = logging.GetLogger()


func Create(ctx context.Context, path string) (*instanceDomain.Instance, error) {
	return instanceDomain.New(ctx, path)
}


func GetByID(ctx context.Context, id string) (*instanceDomain.Instance, error) {
	return instanceDomain.GetByID(ctx, id)
}
