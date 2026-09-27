# BlockchainStore

`blockchainstore` is a Go library for persisting and querying blockchain blocks in SQL databases (e.g., SQLite, MySQL, PostgreSQL). It provides a flexible and simple interface for managing block objects and block chains.

## Features

- **Database Agnostic SQL Store**: Uses `dracory/neat` for SQL query generation across various database drivers.
- **Auto Migration**: Automatically creates and manages the block store database table schema.
- **Rich Querying**: Flexible block listing options including filtering by ID, limit, offset, sorting, and soft deletion.
- **Clean Block Abstraction**: Methods for reading and writing block header fields (`PreviousHash`, `ThisHash`, `Data`, `Timestamp`, etc.).

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

	// Create a new block
	block := blockchainstore.NewBlock()
	block.SetData("Hello Blockchain")
	block.SetPreviousHash("0000000000000000000000000000000000000000000000000000000000000000")
	block.SetThisHash("e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855")

	if err := store.BlockCreate(ctx, block); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Block Created! ID: %s\n", block.ID())

	// Retrieve the block
	found, err := store.BlockFindByID(ctx, block.ID())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Retrieved Block Data: %s\n", found.Data())
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
