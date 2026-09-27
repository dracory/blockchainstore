package main

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/dracory/blockchainstore"
	_ "modernc.org/sqlite"
)

// calculateHash generates a SHA256 hash for block data + previous hash
func calculateHash(data string, prevHash string) string {
	record := data + prevHash
	h := sha256.New()
	h.Write([]byte(record))
	return hex.EncodeToString(h.Sum(nil))
}

func runAdvancedExample(db *sql.DB) error {
	ctx := context.Background()

	// 1. Initialize store with options enabled
	store, err := blockchainstore.NewStore(blockchainstore.NewStoreOptions{
		BlockTableName:     "advanced_blocks",
		DB:                 db,
		AutomigrateEnabled: true,
		DebugEnabled:       false,
	})
	if err != nil {
		return fmt.Errorf("failed to create store: %w", err)
	}

	// 2. Build a blockchain sequence of chained blocks
	var prevHash = "0000000000000000000000000000000000000000000000000000000000000000"
	var createdBlocks []*blockchainstore.Block

	for i := 1; i <= 5; i++ {
		data := fmt.Sprintf("Transaction Block #%d", i)
		hash := calculateHash(data, prevHash)

		block := blockchainstore.NewBlock()
		block.SetData(data)
		block.SetPreviousHash(prevHash)
		block.SetThisHash(hash)

		if err := store.BlockCreate(ctx, block); err != nil {
			return fmt.Errorf("failed to create block %d: %w", i, err)
		}

		createdBlocks = append(createdBlocks, block)
		prevHash = hash
	}

	fmt.Printf("Successfully created and chained %d blocks.\n", len(createdBlocks))

	// 3. Query blocks with pagination and sorting
	options := blockchainstore.BlockQueryOptions{
		Limit:     3,
		Offset:    0,
		OrderBy:   "created_at",
		SortOrder: "asc",
	}

	paginatedBlocks, err := store.BlockList(ctx, options)
	if err != nil {
		return fmt.Errorf("failed to list blocks with query options: %w", err)
	}

	fmt.Printf("Fetched %d blocks (Page 1):\n", len(paginatedBlocks))
	for _, b := range paginatedBlocks {
		fmt.Printf(" - ID: %s | Data: %s | Hash: %s\n", b.ID(), b.Data(), b.ThisHash())
	}

	// 4. Update block previous hash
	targetBlock := createdBlocks[1]
	targetBlock.SetPreviousHash("updated_previous_hash_value")

	if err := store.BlockUpdate(ctx, targetBlock); err != nil {
		return fmt.Errorf("failed to update block: %w", err)
	}

	updatedBlock, err := store.BlockFindByID(ctx, targetBlock.ID())
	if err != nil {
		return fmt.Errorf("failed to find block after update: %w", err)
	}
	fmt.Printf("Updated Block ID %s previous hash to: %s\n", updatedBlock.ID(), updatedBlock.PreviousHash())

	// 5. Delete a block by ID
	blockToDelete := createdBlocks[4]
	if err := store.BlockDeleteByID(ctx, blockToDelete.ID()); err != nil {
		return fmt.Errorf("failed to delete block by ID: %w", err)
	}

	deletedBlock, err := store.BlockFindByID(ctx, blockToDelete.ID())
	if err != nil {
		return fmt.Errorf("error checking deleted block: %w", err)
	}
	if deletedBlock != nil {
		return fmt.Errorf("expected block %s to be deleted", blockToDelete.ID())
	}

	fmt.Printf("Successfully deleted block ID: %s\n", blockToDelete.ID())

	return nil
}

func main() {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	if err := runAdvancedExample(db); err != nil {
		log.Fatalf("advanced example failed: %v", err)
	}
}
