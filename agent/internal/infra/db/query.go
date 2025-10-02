package db

import (
	"database/sql"
	"fmt"
	"strings"
)

// Select executes a SELECT query and returns multiple rows as a slice of maps.
// Each map represents a row with column names as keys and their values as interfaces.
//
// Parameters:
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
//   results, err := Select("users",
//       []string{"id", "name"},
//       map[string]any{"status": "active"})
func Select( tableName  string,
             columns    []string,
             conditions map[string]any,
             args       ...any,
            ) ([]map[string]any, error) {
    if err := hasRequiredParams(tableName, columns); err != nil {
        return nil, err
    }

    var finalArgs []any
    sqlStmt := fmt.Sprintf( "SELECT %s FROM %s",
                            strings.Join(columns, ", "),
                            tableName )

    if len(conditions) > 0 {
        wherePairs, whereValues := buildWhereClause(conditions)
        sqlStmt += fmt.Sprintf(" WHERE %s", strings.Join(wherePairs, " AND "))
        finalArgs = whereValues
    }
    if len(args) > 0 {
        finalArgs = append(finalArgs, args...)
    }
    log.Debug(sqlStmt)

    rows, err := db.Query(sqlStmt, finalArgs...)
    if err != nil {
        return nil, fmt.Errorf("query failed: %v", err)
    }
    defer rows.Close()

    return rowsAsMap(rows, columns)
}

// rowsAsMap converts sql.Rows to a slice of maps where each map represents a row.
// The map keys are column names and values are the corresponding row values.
//
// Parameters:
//   - rows: A pointer to sql.Rows containing the query results
//   - columns: Slice of column names in the order they appear in the query
//
// Returns:
//   - []map[string]any: Slice of maps where each map represents a row
//   - error: Non-nil if an error occurs during row scanning or iteration
//
// Example:
//   rows, err := db.Query("SELECT id, name FROM users")
//   if err != nil {
//       return nil, err
//   }
//   defer rows.Close()
//
//   results, err := rowsAsMap(rows, []string{"id", "name"})
//
// Notes:
//   - The caller is responsible for closing the rows
//   - Column names must match the order in the query
//   - Values are returned as interfaces and may need type assertion
func rowsAsMap(rows *sql.Rows, columns []string) ([]map[string]any, error) {
    var results []map[string]any

    tmpValues   := make([]any, len(columns))
    scanTargets := make([]any, len(columns))
    for i := range scanTargets {
        scanTargets[i] = &tmpValues[i]
    }

    for rows.Next() {
        if err := rows.Scan(scanTargets...); err != nil {
            return nil, fmt.Errorf("error scanning row: %w", err)
        }

        row := make(map[string]any, len(columns))
        for i, col := range columns {
            row[col] = tmpValues[i]
        }

        results = append(results, row)
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("error iterating rows: %w", err)
    }

    return results, nil
}


// SelectOne executes a SELECT query expecting exactly one row result.
// It returns an error if zero or multiple rows are found.
//
// Parameters:
//   - tableName: The name of the table to query
//   - columns: Slice of column names to select
//   - conditions: Map of column-value pairs for WHERE clause
//
// Returns:
//   - map[string]any: A map containing the row data with column names as keys
//   - error: Non-nil if no row found, multiple rows found, or other errors occur
//
// Example:
//   result, err := SelectOne("users",
//       []string{"id", "name"},
//       map[string]any{"id": 1})
func SelectOne(tableName string, columns []string, conditions map[string]any) (map[string]any, error) {
    results, err := Select(tableName, columns, conditions)
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

// getColumnValue is a generic function that retrieves a single column value
// and performs type assertion to the specified type T.
//
// Type Parameter T must be one of: string, []byte, int, or float64
//
// Parameters:
//   - tableName: The name of the table to query
//   - column: The name of the column to retrieve
//   - conditions: Map of column-value pairs for WHERE clause
//
// Returns:
//   - T: The column value as the specified type
//   - error: Non-nil if query fails or type assertion fails
func getColumnValue[T string | []byte | int | float64](
    tableName string,
    column string,
    conditions map[string]any,
) (T, error) {
    result, err := SelectOne(tableName, []string{column}, conditions)
    if err != nil {
        return *new(T), err
    }

    value, ok := result[column].(T)
    if !ok {
        return *new(T), fmt.Errorf("type assertion failed for column %s", column)
    }

    return value, nil
}

// GetAsString retrieves a single column value as a string.
//
// Parameters:
//   - tableName: The name of the table to query
//   - column: The name of the column to retrieve
//   - conditions: Map of column-value pairs for WHERE clause
//
// Returns:
//   - string: The column value
//   - error: Non-nil if query fails or value cannot be converted to string
func GetAsString(tableName, column string, conditions map[string]any) (string, error) {
    return getColumnValue[string](tableName, column, conditions)
}

// GetAsBytes retrieves a single column value as a byte slice.
//
// Parameters:
//   - tableName: The name of the table to query
//   - column: The name of the column to retrieve
//   - conditions: Map of column-value pairs for WHERE clause
//
// Returns:
//   - []byte: The column value
//   - error: Non-nil if query fails or value cannot be converted to []byte
func GetAsBytes(tableName, column string, conditions map[string]any) ([]byte, error) {
    return getColumnValue[[]byte](tableName, column, conditions)
}

// GetAsInt retrieves a single column value as an integer.
//
// Parameters:
//   - tableName: The name of the table to query
//   - column: The name of the column to retrieve
//   - conditions: Map of column-value pairs for WHERE clause
//
// Returns:
//   - int: The column value
//   - error: Non-nil if query fails or value cannot be converted to int
func GetAsInt(tableName, column string, conditions map[string]any) (int, error) {
    return getColumnValue[int](tableName, column, conditions)
}

// GetAsFloat retrieves a single column value as a float64.
//
// Parameters:
//   - tableName: The name of the table to query
//   - column: The name of the column to retrieve
//   - conditions: Map of column-value pairs for WHERE clause
//
// Returns:
//   - float64: The column value
//   - error: Non-nil if query fails or value cannot be converted to float64
func GetAsFloat(tableName, column string, conditions map[string]any) (float64, error) {
    return getColumnValue[float64](tableName, column, conditions)
}
