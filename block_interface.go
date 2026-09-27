package blockchainstore

// BlockInterface defines the interface for blockchain block operations with fluent setters.
type BlockInterface interface {
	// GetID returns the unique identifier of the block.
	GetID() string
	// ID returns the unique identifier of the block.
	ID() string
	// SetID sets the unique identifier of the block.
	SetID(id string) BlockInterface

	// GetTimestamp returns the creation timestamp of the block.
	GetTimestamp() string
	// Timestamp returns the creation timestamp of the block.
	Timestamp() string
	// SetTimestamp sets the creation timestamp of the block.
	SetTimestamp(timestamp string) BlockInterface

	// GetPreviousHash returns the hash of the previous block in the chain.
	GetPreviousHash() string
	// PreviousHash returns the hash of the previous block in the chain.
	PreviousHash() string
	// SetPreviousHash sets the hash of the previous block in the chain.
	SetPreviousHash(previousHash string) BlockInterface

	// GetThisHash returns the hash of the current block.
	GetThisHash() string
	// ThisHash returns the hash of the current block.
	ThisHash() string
	// SetThisHash sets the hash of the current block.
	SetThisHash(thisHash string) BlockInterface

	// Data returns the payload data of the block.
	Data() string
	// SetData sets the payload data of the block.
	SetData(data string) BlockInterface

	// CalculateHash returns the computed SHA-256 hash for this block based on its data and previous hash.
	CalculateHash() string
	// CalculateThisHash computes the block hash and sets it on this block.
	CalculateThisHash() BlockInterface

	// DataObject methods
	GetDataMap() map[string]string
	GetDataChanged() map[string]string
	MarkAsNotDirty(columns ...string)
	MarkAsDirty(columns ...string)
	Get(key string) string
	Set(key string, value string)
	Hydrate(data map[string]string)
	IsDirty() bool
}
