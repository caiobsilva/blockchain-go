package main

import (
	"crypto/sha256" //creates sha256 hashes
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type Block struct {
	Index     int
	Timestamp string
	Data      string
	Hash      string
	PrevHash  string
}

var Blockchain []Block

func calculateHash(block Block) string {
	record := string(block.Index) + block.Timestamp + block.Data + block.PrevHash
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)

	return hex.EncodeToString(hashed)
}

func generateBlock(oldBlock Block, data string) (Block, error) {
	var newBlock Block

	t := time.Now()

	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.Data = data
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Hash = calculateHash(newBlock)

	return newBlock, nil
}

//checking if blocks haven't been tampared with: We do this by checking Index to make sure they've incremented as expected. We also check to make sure our PrevHash is indeed the same as the hash of the previous block. Lastly, we want to double check the hash of the current block by running the calculateHash function again on the current block:
func isBlockValid(newBlock, oldBlock Block) bool {
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}

	if oldBlock.Hash != newBlock.PrevHash {
		return false
	}

	if calculateHash(newBlock) != newBlock.Hash {
		return false
	}

	return true
}

//this is a classic blockchain issue: What if we run into an issue where two nodes of our blockchain ecosystem both added blocks to their chains and we received them both. Which one do we pick as the source of truth? The longest chain.
func replaceChain(newBlocks []Block) {
	if len(newBlocks) > len(Blockchain) {
		Blockchain = newBlocks
	}
}
