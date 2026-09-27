package blockchainstore

import (
	"database/sql"
	"errors"

	"github.com/dracory/neat"
)

// NewStoreOptions define the options for creating a new block store
type NewStoreOptions struct {
	BlockTableName     string
	DB                 *sql.DB
	DbDriverName       string
	TimeoutSeconds     int64
	AutomigrateEnabled bool
	DebugEnabled       bool
}

// NewStore creates a new block store
func NewStore(opts NewStoreOptions) (*Store, error) {
	if opts.BlockTableName == "" {
		return nil, errors.New("block store: BlockTableName is required")
	}

	if opts.DB == nil {
		return nil, errors.New("block store: DB is required")
	}

	neatDB, err := neat.NewFromSQLDB(opts.DB)
	if err != nil {
		return nil, err
	}

	store := &Store{
		blockTableName:     opts.BlockTableName,
		automigrateEnabled: opts.AutomigrateEnabled,
		db:                 neatDB,
		debugEnabled:       opts.DebugEnabled,
	}

	store.timeoutSeconds = 2 * 60 * 60 // 2 hours

	if store.automigrateEnabled {
		if err := store.AutoMigrate(); err != nil {
			return nil, err
		}
	}

	return store, nil
}
