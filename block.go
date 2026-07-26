package blockchainstore

import (
	"github.com/dracory/dataobject"
	"github.com/dracory/uid"
	"github.com/dromara/carbon/v2"
)

type Block struct {
	dataobject.DataObject
}

func NewBlock() *Block {
	block := &Block{}
	block.SetID(uid.HumanUid())
	block.SetTimestamp(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	block.Set("previous_hash", "") // the hash of the previous block
	block.Set("this_hash", "")     // the hash of the current block
	block.Set("data", "")          // the data or transactions (body info)

	return block
}

func NewBlockFromExistingData(data map[string]string) *Block {
	block := &Block{}
	block.Hydrate(data)
	return block
}

func NewBlockFromJSON(json string) *Block {
	data, err := fromJSON(json, nil)
	if err != nil {
		return nil
	}
	blockMap := anyToMapStringString(data)
	return NewBlockFromExistingData(blockMap)
}

// == GETTERS and SETTERS =====================================================

// Timestamp is the time when the block was created
//
// Returns:
// - string: the time when the block was created
func (o *Block) Timestamp() string {
	return o.Get("timestamp")
}

func (o *Block) SetTimestamp(timestamp string) *Block {
	o.Set("timestamp", timestamp)
	return o
}

func (o *Block) PreviousHash() string {
	return o.Get("previous_hash")
}

func (o *Block) SetPreviousHash(previousHash string) *Block {
	o.Set("previous_hash", previousHash)
	return o
}

func (o *Block) ThisHash() string {
	return o.Get("this_hash")
}

func (o *Block) SetThisHash(thisHash string) *Block {
	o.Set("this_hash", thisHash)
	return o
}

func (o *Block) Data() string {
	return o.Get("data")
}

func (o *Block) SetData(data string) *Block {
	o.Set("data", data)
	return o
}
