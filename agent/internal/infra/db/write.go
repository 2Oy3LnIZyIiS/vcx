package db

import (
	"database/sql"
	"fmt"
	"strings"

	"vcx/agent/internal/infra/db/consts"
)


func execute(sqlStmt string, args []any) (sql.Result, error) {
    log.Debug("Executing SQL statement", "statement", sqlStmt)

    result, err := db.Exec(sqlStmt, args...)
    if err != nil {
        log.Error( "Execution Failed",
                   "ERROR", err,
                   "SQL", sqlStmt,
                   "args", args )
        return nil, fmt.Errorf("execution failed: %w", err)
    }
    return result, nil
}


func executeAndGetRowsAffected(sqlStmt string, args []any) (int64, error) {
    result, err := execute(sqlStmt, args)
    if err != nil {
        return 0, err
    }
    return result.RowsAffected()
}


func CreateTable(tableName string, schema map[string]string) bool {
    if tableExists(tableName) {
        return false
    }
    schema[consts.ID]           = "TEXT PRIMARY KEY"
    schema[consts.CREATIONDATE] = "TEXT NOT NULL"
    schema[consts.LMD]          = "TEXT NOT NULL"
    schema[consts.LMU]          = "TEXT NOT NULL"
    schema[consts.GUID]         = "TEXT"
    _, err := execute( fmt.Sprintf( "CREATE TABLE IF NOT EXISTS %s (%s);",
                          tableName,
                          schemaToString(schema)),
             nil )
    if err != nil {
        // TODO: handle errors at this point
        // ... should errors panic?
        return false
    }
    log.Debug("Table created", "table", tableName)

    return true
}


func Insert( tableName string,
             data      map[string]any,
             schema    map[string]string,
           ) (string, error) {
    if err := hasRequiredParams(tableName, data); err != nil {
        return "", err }
    addLMULMD(data)
    if schema != nil{
        // Set default values
        setDefaults(schema, data) }
    columns, placeholders, values := extractColumnsAndValues(data)
    sqlStmt := fmt.Sprintf( "INSERT INTO %s (%s) VALUES (%s)",
                            tableName,
                            strings.Join(columns, ", "),
                            strings.Join(placeholders, ", "))
    if _, err := execute(sqlStmt, values);  err != nil {
        return "", fmt.Errorf("error inserting data: %s", err) }
    return data[consts.ID].(string), nil
}


func Update( tableName string,
             data map[string]any,
             conditions map[string]any,
           ) (int64, error) {
    if err := hasRequiredParams(tableName, data, conditions); err != nil {
        return 0, err }
    addLMULMD(data)
    columns, _, values      := extractColumnsAndValues(data)
    setPairs                := buildSetClause(columns, nil)
    wherePairs, whereValues := buildWhereClause(conditions)
    // values                   = append(values, whereValues...)
    // TODO: check if the above append wherevalues is in original code
    sqlStmt                 := fmt.Sprintf( "UPDATE %s SET %s",
                                            tableName,
                                            strings.Join(setPairs, ", ") )
    if len(conditions) > 0 {
        sqlStmt += fmt.Sprintf(" WHERE %s", strings.Join(wherePairs, " AND "))
        values = append(values, whereValues) }
    return executeAndGetRowsAffected(sqlStmt, values)
}


func Upsert( tableName string,
             data map[string]any,
             // keyColumns []string,
             // columnsIgnoredOnUpdate []string,
             schema map[string]string,
           ) (string, error) {
    // if err := hasRequiredParams(tableName, data, keyColumns); err != nil {
    //     return nil, err
    // }
    if err := hasRequiredParams(tableName, data); err != nil {
        return "", err
    }
    addLMULMD(data)

    if schema != nil {
        setDefaults(schema, data)
    }

    columns, placeholders, values := extractColumnsAndValues(data)

    // Create a map  of columns to exclude
    // excludeColumns := make(map[string]bool, len(keyColumns) + len(columnsIgnoredOnUpdate))
    excludeColumns := make(map[string]bool)
    excludeColumns[consts.ID] = true
    for key, _ := range schema {
        if _, exists := defaultValueGenerators[key]; exists {
            excludeColumns[key] = true
        }
    }
    // for _, col := range keyColumns {
    //     excludeColumns[col] = slices.Contains(keyColumns, col)
    // }
    // for _, col := range columns {
    //     excludeColumns[col] = slices.Contains(columnsIgnoredOnUpdate, col)
    // }
    setPairs := buildSetClause(columns, excludeColumns)

    // add values for update to value list
    for col, val := range data {
        if ! excludeColumns[col] {
            values = append(values, val)
        }
    }

    // Construct the full upsert statement
    sqlStmt := fmt.Sprintf(
        "INSERT INTO %s (%s) VALUES (%s) ON CONFLICT(%s) DO UPDATE SET %s",
        tableName,
        strings.Join(columns,      ", "),
        strings.Join(placeholders, ", "),
        // strings.Join(keyColumns,   ", "),
        consts.ID,
        strings.Join(setPairs,     ", "),
    )
    // values = append(values, updateValues)
    _, error := execute(sqlStmt, values)
    if error != nil {
        return "", error
    }
    return schema[consts.ID], nil
}
