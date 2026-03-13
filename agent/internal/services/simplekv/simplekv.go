package simplekv

import (
	"context"

	"vcx/agent/internal/infra/db/store/simplekv"
	"vcx/pkg/logging"
)

var log = logging.GetLogger()

// GetString retrieves a string value by key.
func GetString(ctx context.Context, key string) (string, error) {
	value, err := simplekv.GetAsString(ctx, key)
	if err != nil {
		log.Debug("SimpleKV get failed", "key", key, "error", err)
		return "", err
	}
	return value, nil
}

// SetString stores a string value by key.
func SetString(ctx context.Context, key, value string) error {
	err := simplekv.SetString(ctx, key, value)
	if err != nil {
		log.Error("SimpleKV set failed", "key", key, "error", err)
		return err
	}
	log.Debug("SimpleKV set success", "key", key)
	return nil
}

// GetInt retrieves an integer value by key.
func GetInt(ctx context.Context, key string) (int, error) {
	value, err := simplekv.GetAsInt(ctx, key)
	if err != nil {
		log.Debug("SimpleKV get failed", "key", key, "error", err)
		return 0, err
	}
	return value, nil
}

// SetInt stores an integer value by key.
func SetInt(ctx context.Context, key string, value int) error {
	err := simplekv.SetInt(ctx, key, value)
	if err != nil {
		log.Error("SimpleKV set failed", "key", key, "error", err)
		return err
	}
	log.Debug("SimpleKV set success", "key", key)
	return nil
}

// GetBytes retrieves a byte slice value by key.
func GetBytes(ctx context.Context, key string) ([]byte, error) {
	value, err := simplekv.GetAsBytes(ctx, key)
	if err != nil {
		log.Debug("SimpleKV get failed", "key", key, "error", err)
		return nil, err
	}
	return value, nil
}

// SetBytes stores a byte slice value by key.
func SetBytes(ctx context.Context, key string, value []byte) error {
	err := simplekv.SetBytes(ctx, key, value)
	if err != nil {
		log.Error("SimpleKV set failed", "key", key, "error", err)
		return err
	}
	log.Debug("SimpleKV set success", "key", key)
	return nil
}