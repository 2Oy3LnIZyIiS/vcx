package simplekv

import (
	"context"
	"fmt"

	"vcx/agent/internal/infra/db"
	"vcx/agent/internal/infra/db/consts"
	// "vcx/pkg/logging"
)

// var log = logging.GetLogger()

const tableName = "simplekv"

const /**Columns*/ (
    COL_KEY   = "key"
    COL_VALUE = consts.VALUE
)

var schema = map[string]string{
    COL_KEY:    "TEXT",
	COL_VALUE: "BLOB NOT NULL",
}


func CreateTable() {
    db.CreateTable(tableName, schema)
}


func GetAsString(ctx context.Context, key string) (string, error) {
	return db.GetAsString(tableName, COL_VALUE, map[string]any{COL_KEY: key})
}


func GetAsBytes(ctx context.Context, key string) ([]byte, error) {
	return db.GetAsBytes(tableName, COL_VALUE, map[string]any{COL_KEY: key})
}


func GetAsInt(ctx context.Context, key string) (int, error) {
	return db.GetAsInt(tableName, COL_VALUE, map[string]any{COL_KEY: key})
}


func GetAsFloat(ctx context.Context, key string) (float64, error) {
	return db.GetAsFloat(tableName, COL_VALUE, map[string]any{COL_KEY: key})
}


func SetValue(ctx context.Context, key string, value any) error {
	if key == "" {
		return fmt.Errorf("key cannot be empty")
	}
	data := map[string]any{
		COL_KEY:   key,
		COL_VALUE: value,
	}
	_, err := db.Upsert(tableName, data, schema)
	return err
}


func SetString(ctx context.Context, key, value string) error {
	return SetValue(ctx, key, value)
}


func SetBytes(ctx context.Context, key string, value []byte) error {
	return SetValue(ctx, key, value)
}


func SetInt(ctx context.Context, key string, value int) error {
	return SetValue(ctx, key, value)
}


func SetFloat(ctx context.Context, key string, value float64) error {
	return SetValue(ctx, key, value)
}
