package db

import (
	"context"
	"fmt"
	"strings"
)

// SelectWithContext executes a SELECT query and returns multiple rows as a slice of maps.
// Each map represents a row with column names as keys and their values as interfaces.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - tableName: The name of the table to query
//   - columns: Slice of column names to select
//   - conditions: Map of column-value pairs for WHERE clause (optional)
//   - args: Additional arguments for the query (optional)
//
// Returns:
//   - []map[string]any: Slice of maps, each containing a row's data
//   - error: Non-nil if any error occurs during query execution or row scanning
//
// Example:
//   ctx := context.Background()
//   results, err := SelectWithContext(ctx, "users",
//       []string{"id", "name"},
//       map[string]any{"status": "active"})
func SelectWithContext(ctx context.Context, tableName string, columns []string, conditions map[string]any, args ...any) ([]map[string]any, error) {
	if err := hasRequiredParams(tableName, columns); err != nil {
		return nil, err
	}

	var finalArgs []any
	sqlStmt := fmt.Sprintf("SELECT %s FROM %s",
		strings.Join(columns, ", "),
		tableName)

	if len(conditions) > 0 {
		wherePairs, whereValues := buildWhereClause(conditions)
		sqlStmt += fmt.Sprintf(" WHERE %s", strings.Join(wherePairs, " AND "))
		finalArgs = whereValues
	}
	if len(args) > 0 {
		finalArgs = append(finalArgs, args...)
	}
	log.Debug(sqlStmt, "args", fmt.Sprintf("%v", finalArgs))

	rows, err := db.QueryContext(ctx, sqlStmt, finalArgs...)
	if err != nil {
		return nil, fmt.Errorf("query failed: %v", err)
	}
	defer rows.Close()

	return rowsAsMap(rows, columns)
}

// SelectOneWithContext executes a SELECT query expecting exactly one row result.
// It returns an error if zero or multiple rows are found.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - tableName: The name of the table to query
//   - columns: Slice of column names to select
//   - conditions: Map of column-value pairs for WHERE clause
//
// Returns:
//   - map[string]any: A map containing the row data with column names as keys
//   - error: Non-nil if no row found, multiple rows found, or other errors occur
//
// Example:
//   ctx := context.Background()
//   result, err := SelectOneWithContext(ctx, "users",
//       []string{"id", "name"},
//       map[string]any{"id": 1})
func SelectOneWithContext(ctx context.Context, tableName string, columns []string, conditions map[string]any) (map[string]any, error) {
	results, err := SelectWithContext(ctx, tableName, columns, conditions)
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no rows found")
	}

	if len(results) > 1 {
		return nil, fmt.Errorf("multiple rows found, expected one")
	}

	return results[0], nil
}

// getColumnValueWithContext is a generic function that retrieves a single column value
// and performs type assertion to the specified type T.
//
// Type Parameter T must be one of: string, []byte, int, or float64
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - tableName: The name of the table to query
//   - column: The name of the column to retrieve
//   - conditions: Map of column-value pairs for WHERE clause
//
// Returns:
//   - T: The column value as the specified type
//   - error: Non-nil if query fails or type assertion fails
func getColumnValueWithContext[T string | []byte | int | float64](
	ctx context.Context,
	tableName string,
	column string,
	conditions map[string]any,
) (T, error) {
	result, err := SelectOneWithContext(ctx, tableName, []string{column}, conditions)
	if err != nil {
		return *new(T), err
	}

	value, ok := result[column].(T)
	if !ok {
		return *new(T), fmt.Errorf("type assertion failed for column %s", column)
	}

	return value, nil
}

// GetAsStringWithContext retrieves a single column value as a string.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - tableName: The name of the table to query
//   - column: The name of the column to retrieve
//   - conditions: Map of column-value pairs for WHERE clause
//
// Returns:
//   - string: The column value
//   - error: Non-nil if query fails or value cannot be converted to string
func GetAsStringWithContext(ctx context.Context, tableName, column string, conditions map[string]any) (string, error) {
	return getColumnValueWithContext[string](ctx, tableName, column, conditions)
}

// GetAsBytesWithContext retrieves a single column value as a byte slice.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - tableName: The name of the table to query
//   - column: The name of the column to retrieve
//   - conditions: Map of column-value pairs for WHERE clause
//
// Returns:
//   - []byte: The column value
//   - error: Non-nil if query fails or value cannot be converted to []byte
func GetAsBytesWithContext(ctx context.Context, tableName, column string, conditions map[string]any) ([]byte, error) {
	return getColumnValueWithContext[[]byte](ctx, tableName, column, conditions)
}

// GetAsIntWithContext retrieves a single column value as an integer.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - tableName: The name of the table to query
//   - column: The name of the column to retrieve
//   - conditions: Map of column-value pairs for WHERE clause
//
// Returns:
//   - int: The column value
//   - error: Non-nil if query fails or value cannot be converted to int
func GetAsIntWithContext(ctx context.Context, tableName, column string, conditions map[string]any) (int, error) {
	return getColumnValueWithContext[int](ctx, tableName, column, conditions)
}

// GetAsFloatWithContext retrieves a single column value as a float64.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - tableName: The name of the table to query
//   - column: The name of the column to retrieve
//   - conditions: Map of column-value pairs for WHERE clause
//
// Returns:
//   - float64: The column value
//   - error: Non-nil if query fails or value cannot be converted to float64
func GetAsFloatWithContext(ctx context.Context, tableName, column string, conditions map[string]any) (float64, error) {
	return getColumnValueWithContext[float64](ctx, tableName, column, conditions)
}
