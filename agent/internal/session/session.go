package session

import (
	"context"
	"fmt"
)

type contextKey string

const (
	ACCOUNTIDKEY contextKey = "accountID"
	CHANGEIDKEY  contextKey = "changeID"
	PROJECTIDKEY contextKey = "projectID"
	BRANCHIDKEY  contextKey = "branchID"
)


func WithAccountID(ctx context.Context, accountID string) context.Context {
	return context.WithValue(ctx, ACCOUNTIDKEY, accountID)
}


func HasAccountID(ctx context.Context) (string, error) {
	if id, ok := ctx.Value(ACCOUNTIDKEY).(string); ok && id != "" {
		return id, nil
	}
	return "", fmt.Errorf("accountID not found in context")
}


func GetAccountID(ctx context.Context) string {
    id, err := HasAccountID(ctx)
    if err != nil {
        panic("accountID not in context - middleware not configured")
    }
    return id
}


func WithChangeID(ctx context.Context, changeID string) context.Context {
	return context.WithValue(ctx, CHANGEIDKEY, changeID)
}


func HasChangeID(ctx context.Context) (string, error) {
	if id, ok := ctx.Value(CHANGEIDKEY).(string); ok && id != "" {
		return id, nil
	}
	return "", fmt.Errorf("changeID not found in context")
}


func GetChangeID(ctx context.Context) string {
    id, err := HasChangeID(ctx)
    if err != nil {
        panic("changeID not in context - middleware not configured")
    }
    return id
}


func WithProjectID(ctx context.Context, projectID string) context.Context {
	return context.WithValue(ctx, PROJECTIDKEY, projectID)
}


func HasProjectID(ctx context.Context) (string, error) {
	if id, ok := ctx.Value(PROJECTIDKEY).(string); ok && id != "" {
		return id, nil
	}
	return "", fmt.Errorf("projectID not found in context")
}


func GetProjectID(ctx context.Context) string {
    id, err := HasProjectID(ctx)
    if err != nil {
        panic("projectID not in context - middleware not configured")
    }
    return id
}


func WithBranchID(ctx context.Context, branchID string) context.Context {
	return context.WithValue(ctx, BRANCHIDKEY, branchID)
}


func HasBranchID(ctx context.Context) (string, error) {
	if id, ok := ctx.Value(BRANCHIDKEY).(string); ok && id != "" {
		return id, nil
	}
	return "", fmt.Errorf("branchID not found in context")
}


func GetBranchID(ctx context.Context) string {
    id, err := HasBranchID(ctx)
    if err != nil {
        panic("branchID not in context - middleware not configured")
    }
    return id
}
