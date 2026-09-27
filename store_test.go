package blockchainstore

import (
	"context"
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"
)

func TestStore(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	opts := NewStoreOptions{
		BlockTableName:     "test_blocks",
		DB:                 db,
		AutomigrateEnabled: true,
	}

	store, err := NewStore(opts)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}

	ctx := context.Background()

	// 1. Create a block using fluent setters
	block := NewBlock().
		SetData("hello blockchain").
		SetPreviousHash("prev_hash_123").
		SetThisHash("this_hash_456")

	err = store.BlockCreate(ctx, block)
	if err != nil {
		t.Fatalf("failed to create block: %v", err)
	}

	if block.GetID() == "" {
		t.Errorf("expected block ID to be non-empty")
	}

	// 2. Find block by ID
	found, err := store.BlockFindByID(ctx, block.GetID())
	if err != nil {
		t.Fatalf("failed to find block by ID: %v", err)
	}
	if found == nil {
		t.Fatalf("expected block to be found")
	}
	if found.Data() != "hello blockchain" {
		t.Errorf("expected data to be 'hello blockchain', got '%s'", found.Data())
	}

	// 3. Update block
	block.SetPreviousHash("updated_prev_hash")
	err = store.BlockUpdate(ctx, block)
	if err != nil {
		t.Fatalf("failed to update block: %v", err)
	}

	found2, err := store.BlockFindByID(ctx, block.GetID())
	if err != nil {
		t.Fatalf("failed to find block after update: %v", err)
	}
	if found2.GetPreviousHash() != "updated_prev_hash" {
		t.Errorf("expected updated previous hash, got '%s'", found2.GetPreviousHash())
	}

	// 4. List blocks
	blocks, err := store.BlockList(ctx, BlockQueryOptions{})
	if err != nil {
		t.Fatalf("failed to list blocks: %v", err)
	}
	if len(blocks) != 1 {
		t.Errorf("expected 1 block in list, got %d", len(blocks))
	}

	// 5. Delete block
	err = store.BlockDelete(ctx, block)
	if err != nil {
		t.Fatalf("failed to delete block: %v", err)
	}

	found3, err := store.BlockFindByID(ctx, block.GetID())
	if err != nil {
		t.Fatalf("failed to find block after delete: %v", err)
	}
	if found3 != nil {
		t.Errorf("expected block to be deleted and not found")
	}

	// 6. Test MigrateDown & MigrateUp
	err = store.MigrateDown(ctx)
	if err != nil {
		t.Fatalf("failed MigrateDown: %v", err)
	}
	err = store.MigrateUp(ctx)
	if err != nil {
		t.Fatalf("failed MigrateUp: %v", err)
	}

	// 7. Test Genesis and Next Block helpers
	gBlock := NewGenesisBlock("Genesis Payload")
	if gBlock.GetPreviousHash() != INITIAL_HASH {
		t.Errorf("expected genesis previous hash to be %s, got %s", INITIAL_HASH, gBlock.GetPreviousHash())
	}
	if gBlock.GetThisHash() == "" {
		t.Errorf("expected genesis hash to be calculated")
	}

	nBlock := NewNextBlock(gBlock, "Block 2 Payload")
	if nBlock.GetPreviousHash() != gBlock.GetThisHash() {
		t.Errorf("expected next block previous hash to match genesis hash")
	}
}
