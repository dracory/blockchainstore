package blockchainstore

import (
	"context"
	"database/sql"
)

type StoreInterface interface {
	MigrateDown(ctx context.Context, tx ...*sql.Tx) error
	MigrateUp(ctx context.Context, tx ...*sql.Tx) error
	EnableDebug(debug bool) StoreInterface

	BlockCreate(ctx context.Context, block BlockInterface) error
	BlockDelete(ctx context.Context, block BlockInterface) error
	BlockDeleteByID(ctx context.Context, blockID string) error
	BlockFindByID(ctx context.Context, id string) (BlockInterface, error)
	BlockList(ctx context.Context, options BlockQueryOptions) ([]BlockInterface, error)
	BlockUpdate(ctx context.Context, block BlockInterface) error
}
