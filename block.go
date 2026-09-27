package blockchainstore

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/dracory/dataobject"
	"github.com/dracory/uid"
	"github.com/dromara/carbon/v2"
)

// INITIAL_HASH is the default previous hash for the genesis block in a chain (64 zeros)
const INITIAL_HASH = "0000000000000000000000000000000000000000000000000000000000000000"

type Block struct {
	dataobject.DataObject
}

var _ BlockInterface = (*Block)(nil)

func NewBlock() BlockInterface {
	block := &Block{}
	block.SetID(uid.HumanUid())
	block.SetTimestamp(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	block.Set(COLUMN_PREVIOUS_HASH, "") // the hash of the previous block
	block.Set(COLUMN_THIS_HASH, "")     // the hash of the current block
	block.Set(COLUMN_DATA, "")          // the data or transactions (body info)

	return block
}

func NewBlockFromExistingData(data map[string]string) BlockInterface {
	block := &Block{}
	block.Hydrate(data)
	return block
}

func NewBlockFromJSON(json string) BlockInterface {
	data, err := fromJSON(json, nil)
	if err != nil {
		return nil
	}
	blockMap := anyToMapStringString(data)
	return NewBlockFromExistingData(blockMap)
}

// NewGenesisBlock creates a new genesis block with INITIAL_HASH as the previous hash and automatically calculates its hash.
func NewGenesisBlock(data string) BlockInterface {
	block := NewBlock().
		SetData(data).
		SetPreviousHash(INITIAL_HASH)
	block.CalculateThisHash()
	return block
}

// NewNextBlock creates a new block linked to the provided previous block and automatically calculates its hash.
func NewNextBlock(prevBlock BlockInterface, data string) BlockInterface {
	prevHash := ""
	if prevBlock != nil {
		prevHash = prevBlock.ThisHash()
	}

	block := NewBlock().
		SetData(data).
		SetPreviousHash(prevHash)
	block.CalculateThisHash()
	return block
}

// CalculateHash generates a SHA256 hash for block data + previous hash.
func CalculateHash(data string, prevHash string) string {
	record := data + prevHash
	h := sha256.New()
	h.Write([]byte(record))
	return hex.EncodeToString(h.Sum(nil))
}

// == GETTERS and SETTERS =====================================================

func (o *Block) GetID() string {
	return o.ID()
}

func (o *Block) SetID(id string) BlockInterface {
	o.DataObject.SetID(id)
	return o
}

func (o *Block) Timestamp() string {
	return o.Get(COLUMN_TIMESTAMP)
}

func (o *Block) GetTimestamp() string {
	return o.Timestamp()
}

func (o *Block) SetTimestamp(timestamp string) BlockInterface {
	o.Set(COLUMN_TIMESTAMP, timestamp)
	return o
}

func (o *Block) PreviousHash() string {
	return o.Get(COLUMN_PREVIOUS_HASH)
}

func (o *Block) GetPreviousHash() string {
	return o.PreviousHash()
}

func (o *Block) SetPreviousHash(previousHash string) BlockInterface {
	o.Set(COLUMN_PREVIOUS_HASH, previousHash)
	return o
}

func (o *Block) ThisHash() string {
	return o.Get(COLUMN_THIS_HASH)
}

func (o *Block) GetThisHash() string {
	return o.ThisHash()
}

func (o *Block) SetThisHash(thisHash string) BlockInterface {
	o.Set(COLUMN_THIS_HASH, thisHash)
	return o
}

func (o *Block) Data() string {
	return o.Get(COLUMN_DATA)
}

func (o *Block) SetData(data string) BlockInterface {
	o.Set(COLUMN_DATA, data)
	return o
}

func (o *Block) CalculateHash() string {
	return CalculateHash(o.Data(), o.PreviousHash())
}

func (o *Block) CalculateThisHash() BlockInterface {
	return o.SetThisHash(o.CalculateHash())
}

func (o *Block) GetDataMap() map[string]string {
	return o.DataObject.Data()
}

func (o *Block) GetDataChanged() map[string]string {
	return o.DataObject.DataChanged()
}
