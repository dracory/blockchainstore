package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/dracory/blockchainstore"
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

	// Create a new genesis block using fluent setters and hash helpers
	block := blockchainstore.NewGenesisBlock("Sample Transaction Data")

	err = store.BlockCreate(ctx, block)
	if err != nil {
		return fmt.Errorf("failed to create block: %w", err)
	}

	fmt.Printf("Created block with ID: %s\n", block.GetID())

	// Retrieve block by ID
	foundBlock, err := store.BlockFindByID(ctx, block.GetID())
	if err != nil {
		return fmt.Errorf("failed to find block by ID: %w", err)
	}

	if foundBlock == nil {
		return fmt.Errorf("block not found")
	}

	fmt.Printf("Retrieved Block ID: %s, Data: %s, Hash: %s\n", foundBlock.GetID(), foundBlock.Data(), foundBlock.ThisHash())
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
