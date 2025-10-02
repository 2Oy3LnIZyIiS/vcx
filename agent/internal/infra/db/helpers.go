package db

import (
	"fmt"
	"strings"

	"vcx/agent/internal/infra/db/consts"
	"vcx/pkg/toolkit/timekit"
)

func schemaToString(schema map[string]string) string {
    columns := make([]string, 0, len(schema))
    for key, value := range schema {
        columns = append(columns, fmt.Sprintf("%s %s", key, value))
    }
    return strings.Join(columns, ", ")
}


func extractColumnsAndValues(data map[string]any) ([]string, []string, []any) {
    columns      := make([]string, 0, len(data))
    placeholders := make([]string, 0, len(data))
    values       := make([]any, 0, len(data))
    for col, val := range data {
        columns      = append(columns, col)
        placeholders = append(placeholders, "?")
        values       = append(values, val)
    }
    return columns, placeholders, values
}

// Helper function to build SET clause for updates
func buildSetClause(columns []string, excludeColumns map[string]bool) []string {
    // log.Debug("Building SET clause for columns: %v; excluding columns: %v", columns, excludeColumns)
    updatePairs := make([]string, 0, len(columns))
    for _, col := range columns {
        // Skip excluded columns in the update
        if excludeColumns != nil {
            if excludeColumns[col] {
                updatePairs = append(updatePairs, fmt.Sprintf("%s = %s", col, col))
            } else {
                updatePairs = append(updatePairs, fmt.Sprintf("%s = excluded.%s", col, col))
            }
        } else {
            updatePairs = append(updatePairs, fmt.Sprintf("%s = ?", col))
        }
    }
    return updatePairs
}


func buildWhereClause(conditions map[string]any) ([]string, []any) {
    // Build WHERE clause
    wherePairs := make([]string, 0, len(conditions))
    values     := make([]any,    0, len(conditions))
    for col, val := range conditions {
        wherePairs = append(wherePairs, fmt.Sprintf("%s = ?", col))
        values     = append(values, val)
    }
    return wherePairs, values
}


func hasRequiredParams(params ...any) error {
    for _, param := range params {
        switch v := param.(type) {
        case map[string]any:
            if len(v) == 0 {
                return fmt.Errorf("data map cannot be empty")
            }
        case map[string]string:
            if len(v) == 0 {
                return fmt.Errorf("schema map cannot be empty")
            }
        case []string:
            if len(v) == 0 {
                return fmt.Errorf("uniqueColumns cannot be empty")
            }
        case string:
            if v == "" {
                return fmt.Errorf("string cannot be empty")
            }
        default:
            return fmt.Errorf("unsupported parameter type: %T", v)
        }
    }
    return nil
}


func addLMULMD(data map[string]any) error {
    if data == nil {
        return fmt.Errorf("data map cannot be nil")
    }
    data[consts.LMD] = timekit.GetDateTime()
    data[consts.LMU] = "local"
    return nil
}
