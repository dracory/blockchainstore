# BlockchainStore

`blockchainstore` is a Go library for persisting and querying blockchain blocks in SQL databases (e.g., SQLite, MySQL, PostgreSQL). It provides a flexible and simple interface for managing block objects and block chains.

## Features

- **Database Agnostic SQL Store**: Uses `dracory/neat` for SQL query generation and database schema migrations across various database drivers.
- **Fluent Setters**: `BlockInterface` supports method chaining for concise block construction and modifications.
- **Convenience Hash & Chain Helpers**: Built-in initial genesis hash, SHA-256 hash calculation, `NewGenesisBlock`, and `NewNextBlock` helpers.
- **Migration Management**: `MigrateUp` and `MigrateDown` schema management using `neat`.
- **Rich Querying**: Flexible block listing options including filtering by ID, limit, offset, and sorting.

## Installation

```bash
go get github.com/dracory/blockchainstore
```

## Quick Start

```go
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/dracory/blockchainstore"
	_ "modernc.org/sqlite"
)

func main() {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	store, err := blockchainstore.NewStore(blockchainstore.NewStoreOptions{
		BlockTableName:     "blocks_block",
		DB:                 db,
		AutomigrateEnabled: true,
	})
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	// Create a new genesis block using fluent setters and built-in hash calculation
	block := blockchainstore.NewGenesisBlock("Hello Blockchain")

	if err := store.BlockCreate(ctx, block); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Block Created! ID: %s, Hash: %s\n", block.GetID(), block.ThisHash())

	// Create next block in the chain
	nextBlock := blockchainstore.NewNextBlock(block, "Second Transaction")
	if err := store.BlockCreate(ctx, nextBlock); err != nil {
		log.Fatal(err)
	}

	// Retrieve the block
	found, err := store.BlockFindByID(ctx, nextBlock.GetID())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Retrieved Next Block Data: %s, Previous Hash: %s\n", found.Data(), found.PreviousHash())
}
```

## Examples

Detailed, runnable examples complete with tests are provided in the [`examples`](./examples) directory:

- [**Basic Example**](./examples/basic) (`examples/basic/main.go`): Demonstrates basic store initialization, creating blocks, and fetching blocks by ID. Includes test `main_test.go`.
- [**Advanced Example**](./examples/advanced) (`examples/advanced/main.go`): Demonstrates building a chain of SHA-256 linked blocks, querying blocks with pagination and sorting, updating blocks, and block deletion. Includes test `main_test.go`.

To run all examples and tests across the project:

```bash
go test ./...
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
