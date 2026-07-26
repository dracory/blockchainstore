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

	// 1. Create a block
	block := NewBlock()
	block.SetData("hello blockchain")
	block.SetPreviousHash("prev_hash_123")
	block.SetThisHash("this_hash_456")

	err = store.BlockCreate(ctx, block)
	if err != nil {
		t.Fatalf("failed to create block: %v", err)
	}

	if block.ID() == "" {
		t.Errorf("expected block ID to be non-empty")
	}

	// 2. Find block by ID
	found, err := store.BlockFindByID(ctx, block.ID())
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

	found2, err := store.BlockFindByID(ctx, block.ID())
	if err != nil {
		t.Fatalf("failed to find block after update: %v", err)
	}
	if found2.PreviousHash() != "updated_prev_hash" {
		t.Errorf("expected updated previous hash, got '%s'", found2.PreviousHash())
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

	found3, err := store.BlockFindByID(ctx, block.ID())
	if err != nil {
		t.Fatalf("failed to find block after delete: %v", err)
	}
	if found3 != nil {
		t.Errorf("expected block to be deleted and not found")
	}
}
