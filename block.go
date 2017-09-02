package main

import (
	"bytes"
	"crypto/sha256"
	"strconv"
	"time"
)

// Block is the block of blockchain.
//
// In the Bitcoin specification Timestamp, PrevBlochHash, and Has are
// block headers, which from a seperate data structure, and
// transactions (Data in this case) is a separate data structure.
// So we're mixsing them here for simplicity.
type Block struct {
	// The current timestamp when the block is created
	Timestamp int64
	// The actual valuable information containing in the block
	Data []byte
	// The hash of the previous block.
	PrevBlockHash []byte
	// The hash of this block.
	Hash  []byte
	Nonce int
}

// SetHash take block fields, concatenate them and calculate a SHA-256 hash
// on the concatenated combination.
func (b *Block) SetHash() {
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	headers := bytes.Join(
		[][]byte{b.PrevBlockHash, b.Data, timestamp}, []byte{})
	hash := sha256.Sum256(headers)

	b.Hash = hash[:]
}

// NewBlock constructs a Block and returns it.
func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{
		time.Now().Unix(), []byte(data), prevBlockHash, []byte{}, 0}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

// NewGenesisBlock creates a genesis block.
// In any blockchain, there must be at least one block, and such block,
// the first in the chain, is called genesis block.
func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}
