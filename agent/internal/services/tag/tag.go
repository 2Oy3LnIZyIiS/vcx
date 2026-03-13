package tag

import (
	"context"

	tagDomain "vcx/agent/internal/domains/tag"
	"vcx/pkg/logging"
)

var log = logging.GetLogger()


func CreateSystemProjectTag(ctx context.Context) (*tagDomain.Tag, error) {
	return tagDomain.NewSystemProjectTag(ctx)
}


func CreateSystemBranchTag(ctx context.Context) (*tagDomain.Tag, error) {
	return tagDomain.NewSystemBranchTag(ctx)
}


func CreateSystemFileTag(ctx context.Context, fileID string) (*tagDomain.Tag, error) {
	return tagDomain.NewSystemFileTag(ctx, fileID)
}


func CreateUserTag(ctx context.Context, name, description string) (*tagDomain.Tag, error) {
	return tagDomain.NewUserTag(ctx, name, description)
}


func CreateUserProjectTag(ctx context.Context, name, description string) (*tagDomain.Tag, error) {
	return tagDomain.NewUserProjectTag(ctx, name, description)
}


func CreateUserBranchTag(ctx context.Context, name, description string) (*tagDomain.Tag, error) {
	return tagDomain.NewUserBranchTag(ctx, name, description)
}


func CreateUserFileTag(ctx context.Context, name, description, fileID string) (*tagDomain.Tag, error) {
	return tagDomain.NewUserFileTag(ctx, name, description, fileID)
}
