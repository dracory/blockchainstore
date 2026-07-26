package blockchainstore

import "github.com/dracory/dataobject"

type Blockchain struct {
	dataobject.DataObject
}

func NewBlockchain() *Blockchain {
	return &Blockchain{}
}
