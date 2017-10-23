package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
	"runtime"
)

var maxNonce = math.MaxInt64

// The difficulty of mining.
// In Bitcoin, "target bits" is the block header storing the diffculty
// at which the block was mined. We won't implement a target ajusting
// algorithm, for now, so we can just define the diffculty as a global
// constant.
//
// 24 is an arbitrary number, our goal is to have a target that takes
// less than 256 bits in memory. And we want the difference to be
// significant enough, but not to big, because the bigger the difference
// the more difficult it's to find a proper hash
const targetBits = 24

// ProofOfWork represents a proof-of-work
type ProofOfWork struct {
	block  *Block
	target *big.Int
}

func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))

	pow := &ProofOfWork{b, target}

	return pow
}

func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash,
			pow.block.HashTransactions(),
			IntToHex(pow.block.Timestamp),
			IntToHex(int64(targetBits)),
			IntToHex(int64(nonce)),
		},
		[]byte{},
	)

	return data
}

type mineResult struct {
	found bool
	nonce int
	hash  []byte
}

func (pow *ProofOfWork) Run() (int, []byte) {
	nonce := 0

	fmt.Printf("Mining a new block")

	noOfGoroutines := runtime.NumCPU() / 2
	nonceRange := maxNonce / noOfGoroutines

	done := make(chan struct{})
	resultCh := make(chan mineResult)

	for i := 0; i < noOfGoroutines; i++ {
		go pow.mining(done, resultCh, nonce, nonce+nonceRange)
		nonce += nonceRange
	}

	for result := range resultCh {
		fmt.Printf("\r%x (%d)", result.hash, result.nonce)
		if result.found {
			close(done)
			fmt.Print("\n\n")
			return result.nonce, result.hash
		}
	}
	panic("Impossible")
}

func (pow *ProofOfWork) mining(done chan struct{}, result chan mineResult,
	nonce int, endOfNonce int) {
	var hashInt big.Int
	var hash [32]byte

	for nonce < endOfNonce {
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(pow.target) == -1 {
			select {
			case result <- mineResult{true, nonce, hash[:]}:
				return
			case <-done:
				return
			}
			panic("Impossible")
		} else {
			select {
			case result <- mineResult{false, nonce, hash[:]}:
			case <-done:
				return
			}
			nonce++
		}
	}
}

func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	return hashInt.Cmp(pow.target) == -1
}
