package main

import (
	"fmt"
	"log"

	"github.com/boltdb/bolt"
)

const dbFile = "blockchain.db"
const blocksBucket = "blocks"

type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

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
	tip []byte
	db  *bolt.DB
}

// AddBlock adds a new Block.
func (bc *Blockchain) AddBlock(data string) {
	var lastHash []byte

	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l"))

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	newBlock := NewBlock(data, lastHash)

	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		err := b.Put(newBlock.Hash, newBlock.Serialize())
		if err != nil {
			log.Fatal(err)
		}

		err = b.Put([]byte("l"), newBlock.Hash)
		if err != nil {
			log.Fatal(err)
		}

		bc.tip = newBlock.Hash

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}

// NewBlockchain creates a blockchain with the genesis block
func NewBlockchain() *Blockchain {
	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)

	if err != nil {
		log.Fatal(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		if b == nil {
			fmt.Println("No existing blockchain found. Creating a new one...")
			genesis := NewGenesisBlock()

			b, err := tx.CreateBucket([]byte(blocksBucket))
			if err != nil {
				log.Fatal(err)
			}

			err = b.Put(genesis.Hash, genesis.Serialize())
			if err != nil {
				log.Fatal(err)
			}

			// 'l' -> the hash of the last block in a chain
			err = b.Put([]byte("l"), genesis.Hash)
			if err != nil {
				log.Fatal(err)
			}

			tip = genesis.Hash
		} else {
			tip = b.Get([]byte("l"))
		}
		return nil
	})

	return &Blockchain{tip, db}
}

func (bc *Blockchain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{bc.tip, bc.db}
}

func (i *BlockchainIterator) Next() *Block {
	var block *Block

	err := i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		encodedBlock := b.Get(i.currentHash)
		block = DeserializeBlock(encodedBlock)

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	i.currentHash = block.PrevBlockHash

	return block
}
