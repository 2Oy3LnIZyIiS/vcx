package db

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"vcx/agent/internal/infra/db/consts"
)

func executeWithContext(ctx context.Context, sqlStmt string, args []any) (sql.Result, error) {
	log.Debug("Executing SQL statement", "statement", sqlStmt)

	result, err := db.ExecContext(ctx, sqlStmt, args...)
	if err != nil {
		log.Error( "Execution Failed",
                   "ERROR", err,
                   "SQL",   sqlStmt,
                   "args",  args)
		return nil, fmt.Errorf("execution failed: %w", err)
	}
	return result, nil
}

func executeAndGetRowsAffectedWithContext(ctx context.Context, sqlStmt string, args []any) (int64, error) {
	result, err := executeWithContext(ctx, sqlStmt, args)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func InsertWithContext(ctx context.Context, tableName string, data map[string]any, schema map[string]string) (string, error) {
	if err := hasRequiredParams(tableName, data); err != nil {
		return "", err
	}
    addMetaTo(schema)
	addLMULMD(data)
	if schema != nil {
		// Set default values
		setDefaults(schema, data)
	}
	columns, placeholders, values := extractColumnsAndValues(data)
	sqlStmt := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		tableName,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "))
	if _, err := executeWithContext(ctx, sqlStmt, values); err != nil {
		return "", fmt.Errorf("error inserting data: %s", err)
	}
	return data[consts.ID].(string), nil
}

func UpdateWithContext(ctx context.Context, tableName string, data map[string]any, conditions map[string]any) (int64, error) {
	if err := hasRequiredParams(tableName, data, conditions); err != nil {
		return 0, err
	}
	addLMULMD(data)
	columns, _, values      := extractColumnsAndValues(data)
	setPairs                := buildSetClause(columns, nil)
	wherePairs, whereValues := buildWhereClause(conditions)
	values                   = append(values, whereValues...)
	sqlStmt                 := fmt.Sprintf( "UPDATE %s SET %s",
                                            tableName,
                                            strings.Join(setPairs, ", "))
	if len(conditions) > 0 {
		sqlStmt += fmt.Sprintf(" WHERE %s", strings.Join(wherePairs, " AND "))
		values   = append(values, whereValues)
	}
	return executeAndGetRowsAffectedWithContext(ctx, sqlStmt, values)
}

func UpsertWithContext(ctx context.Context, tableName string, data map[string]any, schema map[string]string) (string, error) {
	if err := hasRequiredParams(tableName, data); err != nil {
		return "", err
	}
	addLMULMD(data)

	if schema != nil {
		setDefaults(schema, data)
	}

	columns, placeholders, values := extractColumnsAndValues(data)

	// Create a map of columns to exclude
	excludeColumns := make(map[string]bool)
	excludeColumns[consts.ID] = true
	for key, _ := range schema {
		if _, exists := defaultValueGenerators[key]; exists {
			excludeColumns[key] = true
		}
	}
	setPairs := buildSetClause(columns, excludeColumns)

	// add values for update to value list
	for col, val := range data {
		if !excludeColumns[col] {
			values = append(values, val)
		}
	}

	// Construct the full upsert statement
	sqlStmt := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s) ON CONFLICT(%s) DO UPDATE SET %s",
		tableName,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "),
		consts.ID,
		strings.Join(setPairs, ", "),
	)
	_, error := executeWithContext(ctx, sqlStmt, values)
	if error != nil {
		return "", error
	}
	return schema[consts.ID], nil
}
