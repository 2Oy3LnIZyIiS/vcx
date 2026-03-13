package account

import (
	"context"

	"vcx/agent/internal/consts/keys"
	accountDomain "vcx/agent/internal/domains/account"
	"vcx/agent/internal/services/simplekv"
	"vcx/pkg/logging"
)


var log = logging.GetLogger()


// GetOrCreateDefaultAccount gets the default account from SimpleKV or creates one if none exists.
func GetOrCreateDefaultAccount(ctx context.Context) (*accountDomain.Account, error) {
	// Try to get existing default account ID from SimpleKV
	accountID, err := simplekv.GetString(ctx, keys.DEFAULT_ACCOUNT)
	if err == nil && accountID != "" {
		// Load existing account
		account, err := accountDomain.GetByID(ctx, accountID)
		if err == nil {
			log.Info("Using existing default account", "id", account.ID, "name", account.Name)
			return account, nil
		}
		log.Warn("Default account ID found but account doesn't exist, creating new one", "accountID", accountID)
	}

	// Create new default account
	log.Info("Creating new default account")
	account, err := accountDomain.New(ctx, "Default User", "user@example.com", "default")
	if err != nil {
		log.Error("Failed to create default account", "error", err)
		return nil, err
	}

	// Store account ID in SimpleKV
	err = simplekv.SetString(ctx, keys.DEFAULT_ACCOUNT, account.ID)
	if err != nil {
		log.Error("Failed to store default account ID", "error", err)
		return nil, err
	}

	log.Info("Created and stored default account", "id", account.ID, "name", account.Name)
	return account, nil
}
