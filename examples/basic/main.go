package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/gouniverse/blockchainstore"
	_ "modernc.org/sqlite"
)

func runBasicExample(db *sql.DB) error {
	ctx := context.Background()

	// Initialize the store with automatic database table migration
	store, err := blockchainstore.NewStore(blockchainstore.NewStoreOptions{
		BlockTableName:     "blocks_block",
		DB:                 db,
		AutomigrateEnabled: true,
	})
	if err != nil {
		return fmt.Errorf("failed to create store: %w", err)
	}

	// Create a new block
	block := blockchainstore.NewBlock()
	block.SetData("Sample Transaction Data")
	block.SetPreviousHash("0000000000000000000000000000000000000000000000000000000000000000")
	block.SetThisHash("a1b2c3d4e5f67890a1b2c3d4e5f67890a1b2c3d4e5f67890a1b2c3d4e5f67890")

	err = store.BlockCreate(ctx, block)
	if err != nil {
		return fmt.Errorf("failed to create block: %w", err)
	}

	fmt.Printf("Created block with ID: %s\n", block.ID())

	// Retrieve block by ID
	foundBlock, err := store.BlockFindByID(ctx, block.ID())
	if err != nil {
		return fmt.Errorf("failed to find block by ID: %w", err)
	}

	if foundBlock == nil {
		return fmt.Errorf("block not found")
	}

	fmt.Printf("Retrieved Block ID: %s, Data: %s, Hash: %s\n", foundBlock.ID(), foundBlock.Data(), foundBlock.ThisHash())
	return nil
}

func main() {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	if err := runBasicExample(db); err != nil {
		log.Fatalf("example failed: %v", err)
	}
}
