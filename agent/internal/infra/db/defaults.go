package db

import (
	"fmt"

	"vcx/agent/internal/infra/db/consts"
	"vcx/pkg/toolkit/timekit"
	"vcx/pkg/toolkit/uuidkit"
)

// default values
var defaultValueGenerators = map[string]func() string{
    consts.ID:           uuidkit.NewUUIDv7AsString,
    consts.CREATIONDATE: timekit.GetDateTime,
    consts.NAME:         uuidkit.NewShortCode,
}


func setDefaults(schema map[string]string, values map[string]any) {
    for column := range schema {
        if _, exists := values[column]; !exists {
            if defaultValue, err := getDefaultValueFor(column); err == nil {
                values[column] = defaultValue
            }
        }
    }
}


func getDefaultValueFor(column string) (string, error) {
    if generator, exists := defaultValueGenerators[column]; exists {
        return generator(), nil
    }
    return "", fmt.Errorf("no default value for column %q found", column)
}
