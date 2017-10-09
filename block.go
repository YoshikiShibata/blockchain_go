package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"log"
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
	Transactions []*Transaction

	// The hash of the previous block.
	PrevBlockHash []byte

	// The hash of this block.
	Hash []byte

	Nonce int
}

// NewBlock constructs a Block and returns it.
func NewBlock(transactions []*Transaction, prevBlockHash []byte) *Block {
	block := &Block{
		time.Now().Unix(), transactions, prevBlockHash, []byte{}, 0}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

// NewGenesisBlock creates a genesis block.
// In any blockchain, there must be at least one block, and such block,
// the first in the chain, is called genesis block.
func NewGenesisBlock(coinbase *Transaction) *Block {
	return NewBlock([]*Transaction{coinbase}, []byte{})
}

// Serialize serializes a block into []byte.
func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	if err := encoder.Encode(b); err != nil {
		log.Fatal(err)
	}

	return result.Bytes()
}

// DeserializeBlock deserializes []byte into a Block.
func DeserializeBlock(d []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(d))
	if err := decoder.Decode(&block); err != nil {
		log.Fatal(err)
	}

	return &block
}

// HashTransactions take hashes of each transaction, concatenate them,
// and get a hash of the concatenated combination
func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte
	var txHash [32]byte

	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.ID)
	}

	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))

	return txHash[:]
}

// SignTransaction signs inputs of a Transaction
func (bc *Blockchain) SignTransaction(tx *Transaction, privKey ecdsa.PrivateKey) {
	prevTXs := make(map[string]Transaction)

	for _, vin := range tx.Vin {
		prevTX, err := bc.FindTransaction(vin.Txid)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}

	tx.Sign(privKey, prevTXs)
}

func (bc *Blockchain) FindTransaction(ID []byte) (Transaction, error) {
	bci := bc.Iterator()

	for {
		block := bci.Next()

		for _, tx := range block.Transactions {
			if bytes.Compare(tx.ID, ID) == 0 {
				return *tx, nil
			}
		}

		if len(block.PrevBlockHash) == 0 {
			return Transaction{}, errors.New("Transaction is not found")
		}
	}
}
