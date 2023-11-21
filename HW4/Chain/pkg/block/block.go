package block

import (
	"Chain/pkg/pro"
	"crypto/sha256"
	"fmt"
	"google.golang.org/protobuf/proto"
)

// Header provides information about the Block.
// Version is the Block's version.
// PreviousHash is the hash of the previous Block.
// MerkleRoot is the hash of all the Block's Transactions.
// DifficultyTarget is the difficulty of achieving a winning Nonce.
// Nonce is a "number only used once" that satisfies the DifficultyTarget.
// Timestamp is when the Block was successfully mined.
type Header struct {
	Version          uint32
	PreviousHash     string
	MerkleRoot       string
	DifficultyTarget string
	Nonce            uint32
	Timestamp        uint32
}

// Block includes a Header and a slice of Transactions.
type Block struct {
	Header       *Header
	Transactions []*Transaction
}

// EncodeHeader returns a pro.Header given a Header.
func EncodeHeader(header *Header) *pro.Header {
	return &pro.Header{
		Version:          header.Version,
		PreviousHash:     header.PreviousHash,
		MerkleRoot:       header.MerkleRoot,
		DifficultyTarget: header.DifficultyTarget,
		Nonce:            header.Nonce,
		Timestamp:        header.Timestamp,
	}
}

// DecodeHeader returns a Header given a pro.Header.
func DecodeHeader(pheader *pro.Header) *Header {
	return &Header{
		Version:          pheader.GetVersion(),
		PreviousHash:     pheader.GetPreviousHash(),
		MerkleRoot:       pheader.GetMerkleRoot(),
		DifficultyTarget: pheader.GetDifficultyTarget(),
		Nonce:            pheader.GetNonce(),
		Timestamp:        pheader.GetTimestamp(),
	}
}

// Task 1-4
// EncodeBlock returns a pro.Block given a Block.
func EncodeBlock(b *Block) *pro.Block {
  //1. initilize an array of pro.Transaction
	var ptxs []*pro.Transaction
  //2. encode each transaction and append it to ptxs
	
  //3. encode the b by consturcting a pro.Block, and return the pointer

  return nil
}

// Task 1-5
// DecodeBlock returns a Block given a pro.Block.
func DecodeBlock(pb *pro.Block) *Block {
  //1. initilize an array of Transaction
	var txs []*Transaction
  //2. decode each transaction and append it to txs
	
  //3. decode the pb by constructing a Block, and return the pointer
	
}

// Task 1-6
// Hash returns the hash of the block (which is done via the header)
func (block *Block) Hash() string {
  //1. initilize a sha256 instance
	
  //2. encode the block
	
  //3. serialize the encoded block
	
  //4. hash the serialized block
	
  //5. return the result

  return ""
}
