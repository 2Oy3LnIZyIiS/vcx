package store

import (
	"context"
	"fmt"
	"vcx/agent/internal/infra/db"
	"vcx/agent/internal/infra/db/consts"
	"vcx/pkg/logging"
)


var log = logging.GetLogger()


func Create( ctx       context.Context,
             tableName string,
             data      map[string]any,
             schema    map[string]string,
           ) (map[string]any, error) {
    log = logging.GetLogger()
	id, err := db.InsertWithContext(ctx, tableName, data, schema)
    if err != nil {
        return nil, err
    }
    log.Debug(fmt.Sprintf("Record Created with id: %s", id))
    return GetByID(ctx, tableName, id)
}


func Update( ctx        context.Context,
             tableName  string,
             id         string,
             data       map[string]any,
           ) (map[string]any, error) {
	_, err := db.UpdateWithContext(ctx, tableName, data, map[string]any{consts.ID: id})
    if err != nil {
        return nil, err
    }
	return GetByID(ctx, tableName, id)
}


func GetByID( ctx        context.Context,
              tableName  string,
              id         string,
            ) (map[string]any, error) {
	result, err := db.SelectOneWithContext(ctx, tableName, []string{"*"}, map[string]any{consts.ID: id})
	if err != nil {
		return nil, err
	}
	return result, nil
}
