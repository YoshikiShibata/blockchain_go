package main

// Blockchain is just a database with certain structure: it's an ordered,
// back-linked list. Which means that blocks are stored in the insertion
// order and that each block is linked to the previous one. This structure
// allows to quickly get the latest block in a chain and to (efficiently)
// get a block by its hash.
//
// In Golang this structure can be implemented by using an array and a
// map: the array would keep ordered hashes(arrays are ordered in Go),
// and the map would keep hash -> block pairs (maps are unordered).
// But for our blockchain prototype we'll just use an array, because
// we don't need to get blocks by their hash for now.
type Blockchain struct {
	blocks []*Block
}

// AddBlock adds a new Block.
func (bc *Blockchain) AddBlock(data string) {
	prevBlock := bc.blocks[len(bc.blocks)-1]
	newBlock := NewBlock(data, prevBlock.Hash)
	bc.blocks = append(bc.blocks, newBlock)
}

// NewBlockchain creates a blockchain with the genesis block
func NewBlockchain() *Blockchain {
	return &Blockchain{[]*Block{NewGenesisBlock()}}
}
